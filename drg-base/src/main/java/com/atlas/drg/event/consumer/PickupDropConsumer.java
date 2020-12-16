package com.atlas.drg.event.consumer;

import com.atlas.drg.command.PickupDropCommand;
import com.atlas.drg.constant.EventConstants;
import com.atlas.drg.processor.DropProcessor;
import com.atlas.drg.processor.TopicDiscoveryProcessor;
import com.atlas.kafka.consumer.SimpleEventHandler;

public class PickupDropConsumer implements SimpleEventHandler<PickupDropCommand> {
   @Override
   public void handle(Long key, PickupDropCommand event) {
      DropProcessor.pickupDrop(event.dropId(), event.characterId());
   }

   @Override
   public Class<PickupDropCommand> getEventClass() {
      return PickupDropCommand.class;
   }

   @Override
   public String getConsumerId() {
      return "Drop Registry";
   }

   @Override
   public String getBootstrapServers() {
      return System.getenv("BOOTSTRAP_SERVERS");
   }

   @Override
   public String getTopic() {
      return TopicDiscoveryProcessor.getTopic(EventConstants.TOPIC_PICKUP_DROP_COMMAND);
   }
}