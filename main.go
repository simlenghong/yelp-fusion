package main

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"time"
	"github.com/simlenghong/yelp-fusion/yelp"
)

const (
	mFolderLibs string = "libs"
	mFolderHtml string = "html"
)

var (
	mYelpAccessTokenTimeExpired time.Time 
	mYelpAccessToken *yelp.AccessToken 
)

//Render Html pages
func RenderPage(pResponseWriter http.ResponseWriter, pHtmlFileName string) {
	log.Println("Processing " + pHtmlFileName + ".html... " )
	lListHtmlFile := []string{ mFolderHtml + "/base.html", 
				mFolderHtml + "/" + pHtmlFileName + ".html" }

	lTemplatePtr, lErr := template.ParseFiles(lListHtmlFile...)
	if lErr != nil {
		log.Println("Error parsing files: ", lErr)
	}

	lErr = lTemplatePtr.Execute(pResponseWriter, nil)
	if lErr != nil {
		log.Println("Error executing template: ", lErr)
	}
}

//Handler for processing files in libs folder
func HandlerFilesLib(pResponseWriter http.ResponseWriter, pHttpRequest *http.Request) {
	lPathLibFile := pHttpRequest.URL.Path[len("/" + mFolderLibs + "/"):]
	
	if len(lPathLibFile) != 0 {
		lFile, lErr := http.Dir(mFolderLibs + "/").Open(lPathLibFile)
		if lErr != nil {
			log.Println("Error opening lib file: ", lErr)
			http.NotFound(pResponseWriter, pHttpRequest)
		} else {
			lReadSeeker := io.ReadSeeker(lFile)
			http.ServeContent(pResponseWriter, pHttpRequest, lPathLibFile, time.Now(), lReadSeeker)
		}
	}
}

//Handler for processing html/index.html
func HandlerPageIndex(pResponseWriter http.ResponseWriter, pHttpRequest *http.Request) {
	RenderPage(pResponseWriter, "index")
}

//Handler for searching business in yelp
func HandlerSearchBusinessInYelp(pResponseWriter http.ResponseWriter, pRequest *http.Request) {
	if IsYelpAccessTokenExpired() {
		SetYelpAccessToken()
	}

	log.Println("Yelp access token expiring on " + mYelpAccessTokenTimeExpired.String()) 

	var lYelpParamSearchBusiness yelp.ParamSearchBusiness 
	lDecoder := json.NewDecoder(pRequest.Body)
	lErr := lDecoder.Decode(&lYelpParamSearchBusiness)
	if lErr != nil {
	        log.Println("Error decoding Json string: ", lErr)
	}
	defer pRequest.Body.Close()

	log.Println("Term: " + lYelpParamSearchBusiness.Term)
	log.Println("Location: " + lYelpParamSearchBusiness.Location)

	lYelpRespSearchBusinessMin, lErr := yelp.SearchBusinessMin(lYelpParamSearchBusiness, mYelpAccessToken)
    	if lErr != nil {
	        log.Println("Error searching businesses in Yelp: ", lErr)
    	}
	log.Println(lYelpRespSearchBusinessMin.Total)

	pResponseWriter.Header().Set("Content-Type", "application/json")
	json.NewEncoder(pResponseWriter).Encode(lYelpRespSearchBusinessMin)
}

func IsYelpAccessTokenExpired() (bool) {
	if time.Now().After(mYelpAccessTokenTimeExpired) {
		log.Println("Yelp access token expired on " + mYelpAccessTokenTimeExpired.String()) 
		return true
	}
	return false
}

func SetYelpAccessToken() {
        log.Println("Retrieving Yelp access token...")
	mYelpAccessTokenTimeExpired = time.Now()
	mYelpAccessToken = new(yelp.AccessToken)
	lErr := yelp.GetAccessToken(mYelpAccessToken)
	if lErr != nil {
		log.Println("Error retrieving Yelp access token: ", lErr)
	}
	lDuration := time.Duration(mYelpAccessToken.ExpiresIn*1000*1000*1000)
	mYelpAccessTokenTimeExpired = mYelpAccessTokenTimeExpired.Add(lDuration)
}

//Main function
func main(){
	SetYelpAccessToken()

        log.Println("Registering handlers...")
	http.HandleFunc("/" + mFolderLibs + "/", HandlerFilesLib)
	http.HandleFunc("/", HandlerPageIndex)
	http.HandleFunc("/SearchBusiness", HandlerSearchBusinessInYelp)
	
        log.Println("Starting webserver...")
	lErrListenAndServe := http.ListenAndServe(":8000", nil)
	if lErrListenAndServe != nil {
	        log.Println("Error starting webserver: ", lErrListenAndServe)
	}
}
