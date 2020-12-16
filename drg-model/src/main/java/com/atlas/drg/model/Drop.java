package com.atlas.drg.model;

import com.atlas.drg.builder.DropBuilder;

public record Drop(int id, int worldId, int channelId, int mapId, int itemId, int quantity, int meso, int type, int x, int y,
                   int ownerId, Integer ownerPartyId, long dropTime, int dropperId, int dropperX, int dropperY, boolean playerDrop,
                   byte mod, DropStatus status) {
   public Drop reserve() {
      return new DropBuilder(this).setStatus(DropStatus.RESERVED).build();
   }
}
