
# yelp-fusion
Sample website calling Yelp Fusion API written in go. The system will first search in MongoDB based on Term and Location. If the record does not exist, it will search in Yelp and then insert the data from Yelp to MongoDB
  
# GETTING STARTED
1. Replace the following variables with your Yelp Fusion id and password in yelp\yelp.go
   * mClientId = "YOUR_YELP_FUSION_CLIENT_ID"
   * mClientSecret = "YOUR_YELP_FUSION_CLIENT_SECRET"
2. Download "gopkg.in/mgo.v2" and "gopkg.in/mgo.v2/bson"
3. Run main.go
4. Open your browser and type in localhost:8000
5. Enter your search criteria and press the search button
6. Results will be shown in a sortable table format.

# TODO
- [x] MongoDB
- [ ] Error Handling
- [ ] Unit test
- [ ] Input validation
- [ ] Config file for DB and Yelp configs
