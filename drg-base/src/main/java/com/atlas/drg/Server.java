package com.atlas.drg;

import java.net.URI;

import com.atlas.drg.event.consumer.KillMonsterConsumer;
import com.atlas.kafka.consumer.SimpleEventConsumerFactory;
import com.atlas.shared.rest.RestServerFactory;
import com.atlas.shared.rest.RestService;
import com.atlas.shared.rest.UriBuilder;

public class Server {
   public static void main(String[] args) {
      SimpleEventConsumerFactory.create(new KillMonsterConsumer());

      URI uri = UriBuilder.host(RestService.DROP_REGISTRY).uri();
      RestServerFactory.create(uri, "com.atlas.drg.rest");
   }
}
