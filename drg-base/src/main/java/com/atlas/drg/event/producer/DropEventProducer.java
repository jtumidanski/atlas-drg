package com.atlas.drg.event.producer;

import com.atlas.drg.constant.EventConstants;
import com.atlas.drg.event.DropEvent;
import com.atlas.drg.model.Drop;
import com.atlas.kafka.KafkaProducerFactory;
import org.apache.kafka.clients.producer.Producer;
import org.apache.kafka.clients.producer.ProducerRecord;

public class DropEventProducer {
   private static final Object lock = new Object();

   private static volatile DropEventProducer instance;

   private final Producer<Long, DropEvent> producer;

   public static DropEventProducer getInstance() {
      DropEventProducer result = instance;
      if (result == null) {
         synchronized (lock) {
            result = instance;
            if (result == null) {
               result = new DropEventProducer();
               instance = result;
            }
         }
      }
      return result;
   }

   private DropEventProducer() {
      producer = KafkaProducerFactory.createProducer("Drop Registry", System.getenv("BOOTSTRAP_SERVERS"));
   }

   public void createDrop(int worldId, int channelId, Drop drop) {
      String topic = System.getenv(EventConstants.TOPIC_DROP_EVENT);
      long key = produceKey(worldId, channelId);
      producer.send(new ProducerRecord<>(topic, key,
            new DropEvent(drop.id(), drop.itemId(), drop.quantity(), drop.type(), drop.x(), drop.y(), drop.ownerId(),
                  drop.ownerPartyId(), drop.dropTime(), drop.dropperId(), drop.dropperX(), drop.dropperY(), drop.playerDrop(),
                  drop.mod())));
   }

   protected Long produceKey(int worldId, int channelId) {
      return (long) ((worldId * 1000) + channelId);
   }
}
