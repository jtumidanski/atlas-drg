package com.atlas.drg.rest.processor;

import com.atlas.drg.model.Drop;
import com.atlas.drg.rest.attribute.DropAttributes;
import com.atlas.drg.rest.builder.DropAttributesBuilder;

import builder.ResultObjectBuilder;

public final class ResultObjectFactory {

   public static ResultObjectBuilder create(Drop drop) {
      return new ResultObjectBuilder(DropAttributes.class, drop.id())
            .setAttribute(new DropAttributesBuilder()
                  .setWorldId(drop.worldId())
                  .setChannelId(drop.channelId())
                  .setMapId(drop.mapId())
                  .setItemId(drop.itemId())
                  .setQuantity(drop.quantity())
                  .setMeso(drop.meso())
                  .setDropType(drop.type())
                  .setDropX(drop.x())
                  .setDropY(drop.y())
                  .setOwnerId(drop.ownerId())
                  .setOwnerPartyId(drop.ownerPartyId())
                  .setDropTime(drop.dropTime())
                  .setDropperUniqueId(drop.dropperId())
                  .setDropperX(drop.dropperX())
                  .setDropperY(drop.dropperY())
                  .setPlayerDrop(drop.playerDrop())
                  .setMod(drop.mod())
            );
   }
}
