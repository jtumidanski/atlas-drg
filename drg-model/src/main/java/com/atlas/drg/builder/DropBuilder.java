package com.atlas.drg.builder;

import com.atlas.drg.model.Drop;
import com.atlas.drg.model.DropStatus;

public class DropBuilder {
   private final Integer id;

   private Integer worldId;

   private Integer channelId;

   private Integer mapId;

   private Integer itemId;

   private Integer quantity;

   private Integer meso;

   private Integer type;

   private Integer x;

   private Integer y;

   private Integer ownerId;

   private Integer ownerPartyId;

   private Long dropTime;

   private Integer dropperId;

   private Integer dropperX;

   private Integer dropperY;

   private Boolean playerDrop;

   private Byte mod;

   private DropStatus status;

   public DropBuilder(Integer id) {
      this.id = id;
   }

   public DropBuilder(Drop other) {
      this.id = other.id();
      this.worldId = other.worldId();
      this.channelId = other.channelId();
      this.mapId = other.mapId();
      this.itemId = other.itemId();
      this.quantity = other.quantity();
      this.meso = other.meso();
      this.type = other.type();
      this.x = other.x();
      this.y = other.y();
      this.ownerId = other.ownerId();
      this.ownerPartyId = other.ownerPartyId();
      this.dropTime = other.dropTime();
      this.dropperId = other.dropperId();
      this.dropperX = other.dropperX();
      this.dropperY = other.dropperY();
      this.playerDrop = other.playerDrop();
      this.mod = other.mod();
      this.status = other.status();
   }

   public Drop build() {
      return new Drop(id, worldId, channelId, mapId, itemId, quantity, meso, type, x, y, ownerId, ownerPartyId, dropTime,
            dropperId, dropperX, dropperY, playerDrop, mod, status);
   }

   public DropBuilder setWorldId(Integer worldId) {
      this.worldId = worldId;
      return this;
   }

   public DropBuilder setChannelId(Integer channelId) {
      this.channelId = channelId;
      return this;
   }

   public DropBuilder setMapId(Integer mapId) {
      this.mapId = mapId;
      return this;
   }

   public DropBuilder setItemId(Integer itemId) {
      this.itemId = itemId;
      return this;
   }

   public DropBuilder setQuantity(Integer quantity) {
      this.quantity = quantity;
      return this;
   }

   public DropBuilder setMeso(Integer meso) {
      this.meso = meso;
      return this;
   }

   public DropBuilder setType(Integer type) {
      this.type = type;
      return this;
   }

   public DropBuilder setX(Integer x) {
      this.x = x;
      return this;
   }

   public DropBuilder setY(Integer y) {
      this.y = y;
      return this;
   }

   public DropBuilder setOwnerId(Integer ownerId) {
      this.ownerId = ownerId;
      return this;
   }

   public DropBuilder setOwnerPartyId(Integer ownerPartyId) {
      this.ownerPartyId = ownerPartyId;
      return this;
   }

   public DropBuilder setDropTime(Long dropTime) {
      this.dropTime = dropTime;
      return this;
   }

   public DropBuilder setDropperId(Integer dropperId) {
      this.dropperId = dropperId;
      return this;
   }

   public DropBuilder setDropperX(Integer dropperX) {
      this.dropperX = dropperX;
      return this;
   }

   public DropBuilder setDropperY(Integer dropperY) {
      this.dropperY = dropperY;
      return this;
   }

   public DropBuilder setPlayerDrop(Boolean playerDrop) {
      this.playerDrop = playerDrop;
      return this;
   }

   public DropBuilder setMod(Byte mod) {
      this.mod = mod;
      return this;
   }

   public DropBuilder setStatus(DropStatus status) {
      this.status = status;
      return this;
   }
}
