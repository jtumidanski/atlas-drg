package com.atlas.drg.event.consumer;

import com.atlas.drg.processor.DropProcessor;
import com.atlas.kafka.consumer.SimpleEventHandler;
import com.atlas.morg.rest.constant.EventConstants;
import com.atlas.morg.rest.event.MonsterKilledEvent;

public class KillMonsterConsumer implements SimpleEventHandler<MonsterKilledEvent> {
   @Override
   public void handle(Long key, MonsterKilledEvent event) {
      DropProcessor.createDrops(event.worldId(), event.channelId(), event.mapId(), event.uniqueId(), event.monsterId(), event.x()
            , event.y(), event.killerId());
   }

   @Override
   public Class<MonsterKilledEvent> getEventClass() {
      return MonsterKilledEvent.class;
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
      return System.getenv(EventConstants.TOPIC_MONSTER_KILLED_EVENT);
   }
}
