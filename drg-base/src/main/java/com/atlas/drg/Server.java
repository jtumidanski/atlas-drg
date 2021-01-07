package com.atlas.drg;

import java.net.URI;
import java.util.concurrent.Executors;
import java.util.concurrent.TimeUnit;

import com.atlas.drg.constant.RestConstants;
import com.atlas.drg.event.consumer.CancelDropReservationConsumer;
import com.atlas.drg.event.consumer.KillMonsterConsumer;
import com.atlas.drg.event.consumer.PickupDropConsumer;
import com.atlas.drg.event.consumer.ReserveDropConsumer;
import com.atlas.drg.processor.DropProcessor;
import com.atlas.drg.task.DropExpireTask;
import com.atlas.kafka.consumer.SimpleEventConsumerFactory;
import com.atlas.shared.rest.RestServerFactory;
import com.atlas.shared.rest.UriBuilder;

public class Server {
   public static void main(String[] args) {
      Runtime.getRuntime().addShutdownHook(new Thread(DropProcessor::destroyAll));

      SimpleEventConsumerFactory.create(new KillMonsterConsumer());
      SimpleEventConsumerFactory.create(new ReserveDropConsumer());
      SimpleEventConsumerFactory.create(new PickupDropConsumer());
      SimpleEventConsumerFactory.create(new CancelDropReservationConsumer());

      Executors.newSingleThreadScheduledExecutor().scheduleWithFixedDelay(new DropExpireTask(),
            0,
            ConfigurationRegistry.getInstance().getConfiguration().itemExpireCheck,
            TimeUnit.MILLISECONDS);

      URI uri = UriBuilder.host(RestConstants.SERVICE).uri();
      RestServerFactory.create(uri, "com.atlas.drg.rest");
   }
}
