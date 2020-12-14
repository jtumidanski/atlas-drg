package com.atlas.drg;

import com.atlas.drg.event.consumer.KillMonsterConsumer;
import com.atlas.kafka.consumer.SimpleEventConsumerFactory;

public class Server {
   public static void main(String[] args) {
      SimpleEventConsumerFactory.create(new KillMonsterConsumer());
   }
}
