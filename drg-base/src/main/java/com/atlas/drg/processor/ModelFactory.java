package com.atlas.drg.processor;

import com.atlas.dis.rest.attribute.MonsterDropAttributes;
import com.atlas.drg.model.MonsterDrop;

import rest.DataBody;

public final class ModelFactory {
   private ModelFactory() {
   }

   public static MonsterDrop createMonsterDrop(DataBody<MonsterDropAttributes> body) {
      return new MonsterDrop(body.getAttributes().monsterId(), body.getAttributes().itemId(),
            body.getAttributes().maximumQuantity(), body.getAttributes().minimumQuantity(), body.getAttributes().chance());
   }
}
