package reviews;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpHeaders;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestHeader;
import org.springframework.web.bind.annotation.RestController;
import org.springframework.web.client.RestTemplate;

import java.util.Map;

@RestController
public class ReviewsController {
    @Autowired
    private RestTemplate restTemplate;

    @Autowired
    private RatingsService ratingsService;

    @GetMapping("/reviews/{productID}")
    public String bookReviewsById(@PathVariable("productID") int productId, @RequestHeader HttpHeaders requestHeaders) {
        int starsReviewer1 = -1;
        int starsReviewer2 = -1;

        if (ratings_enabled) {
            Ratings ratingsResponse = getRatings(Integer.toString(productId), requestHeaders);
            if (ratingsResponse != null) {
                if (ratingsResponse.getRatings() != null) {
                    Map<String, Integer> ratings = ratingsResponse.getRatings();
                    if (ratings.containsKey("Reviewer1")){
                        starsReviewer1 = ratings.get("Reviewer1");
                    }
                    if (ratings.containsKey("Reviewer2")){
                        starsReviewer2 = ratings.get("Reviewer2");
                    }
                }
            }
        }

        return getJsonResponse(Integer.toString(productId), starsReviewer1, starsReviewer2);
    }

    @GetMapping("/rest")
    public String rest() {
        return restTemplate.getForObject("http://ratings/health", String.class);
    }

    @GetMapping("/health")
    public String health() {
        return "{\"status\": \"Reviews is healthy\"}";
    }

    private final static Boolean ratings_enabled = Boolean.valueOf(System.getenv("ENABLE_RATINGS"));
    private final static String star_color = System.getenv("STAR_COLOR") == null ? "black" : System.getenv("STAR_COLOR");

    private Ratings getRatings(String productId, HttpHeaders requestHeaders) {
        return ratingsService.getRatings(productId, requestHeaders);
    }

    private String getJsonResponse (String productId, int starsReviewer1, int starsReviewer2) {
        String result = "{";
        result += "\"id\": \"" + productId + "\",";
        result += "\"reviews\": [";

        // reviewer 1:
        result += "{";
        result += "  \"reviewer\": \"Reviewer1\",";
        result += "  \"text\": \"An extremely entertaining play by Shakespeare. The slapstick humour is refreshing!\"";
        if (ratings_enabled) {
            if (starsReviewer1 != -1) {
                result += ", \"rating\": {\"stars\": " + starsReviewer1 + ", \"color\": \"" + star_color + "\"}";
            }
            else {
                result += ", \"rating\": {\"error\": \"Ratings service is currently unavailable\"}";
            }
        }
        result += "},";

        // reviewer 2:
        result += "{";
        result += "  \"reviewer\": \"Reviewer2\",";
        result += "  \"text\": \"Absolutely fun and entertaining. The play lacks thematic depth when compared to other plays by Shakespeare.\"";
        if (ratings_enabled) {
            if (starsReviewer2 != -1) {
                result += ", \"rating\": {\"stars\": " + starsReviewer2 + ", \"color\": \"" + star_color + "\"}";
            }
            else {
                result += ", \"rating\": {\"error\": \"Ratings service is currently unavailable\"}";
            }
        }
        result += "}";

        result += "]";
        result += "}";

        return result;
    }
}
