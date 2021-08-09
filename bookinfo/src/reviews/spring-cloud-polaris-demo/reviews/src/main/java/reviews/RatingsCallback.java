package reviews;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.http.HttpHeaders;
import org.springframework.stereotype.Component;

import java.util.Collections;
import java.util.Random;

@Component
public class RatingsCallback implements RatingsService{
    private static final Logger LOG = LoggerFactory.getLogger(RatingsCallback.class);

    @Override
    public Ratings getRatings(String productID, HttpHeaders requestHeaders) {
        LOG.error("getRatings encountered circuit breaker !!!");
        Random random = new Random();
        Ratings ratings = new Ratings();
        ratings.setId(Integer.parseInt(productID));
        ratings.setRatings(Collections.singletonMap("MockReviewer1", random.nextInt(5)));
        return ratings;
    }
}
