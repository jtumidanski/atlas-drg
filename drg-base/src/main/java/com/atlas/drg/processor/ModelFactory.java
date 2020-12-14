package com.atlas.drg.processor;

import com.atlas.dis.rest.attribute.MonsterDropAttributes;
import com.atlas.drg.model.Monster;
import com.atlas.drg.model.MonsterDrop;
import com.atlas.morg.rest.attribute.MonsterAttributes;

import rest.DataBody;

public final class ModelFactory {
   private ModelFactory() {
   }

   public static MonsterDrop createMonsterDrop(DataBody<MonsterDropAttributes> body) {
      return new MonsterDrop(body.getAttributes().monsterId(), body.getAttributes().itemId(),
            body.getAttributes().maximumQuantity(), body.getAttributes().minimumQuantity(), body.getAttributes().chance());
   }

   public static Monster createMonster(DataBody<MonsterAttributes> body) {
      return new Monster(Integer.parseInt(body.getId()), body.getAttributes().monsterId(),
            body.getAttributes().worldId(), body.getAttributes().channelId(), body.getAttributes().mapId(),
            body.getAttributes().x(), body.getAttributes().y());
   }
}
