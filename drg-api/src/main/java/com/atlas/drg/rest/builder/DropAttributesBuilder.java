package com.atlas.drg.rest.builder;

import com.app.common.builder.RecordBuilder;
import com.atlas.drg.rest.attribute.DropAttributes;

import builder.AttributeResultBuilder;

public class DropAttributesBuilder extends RecordBuilder<DropAttributes, DropAttributesBuilder> implements AttributeResultBuilder {
   private Integer worldId;

   private Integer channelId;

   private Integer mapId;

   private Integer itemId;

   private Integer quantity;

   private Integer meso;

   private Integer dropType;

   private Integer dropX;

   private Integer dropY;

   private Integer ownerId;

   private Integer ownerPartyId;

   private Long dropTime;

   private Integer dropperUniqueId;

   private Integer dropperX;

   private Integer dropperY;

   private Boolean playerDrop;

   private Byte mod;

   @Override
   public DropAttributes construct() {
      return new DropAttributes(worldId, channelId, mapId, itemId, quantity, meso, dropType, dropX, dropY, ownerId, ownerPartyId,
            dropTime, dropperUniqueId, dropperX, dropperY, playerDrop, mod);
   }

   @Override
   public DropAttributesBuilder getThis() {
      return this;
   }

   public DropAttributesBuilder setWorldId(Integer worldId) {
      this.worldId = worldId;
      return getThis();
   }

   public DropAttributesBuilder setChannelId(Integer channelId) {
      this.channelId = channelId;
      return getThis();
   }

   public DropAttributesBuilder setMapId(Integer mapId) {
      this.mapId = mapId;
      return getThis();
   }

   public DropAttributesBuilder setItemId(Integer itemId) {
      this.itemId = itemId;
      return getThis();
   }

   public DropAttributesBuilder setQuantity(Integer quantity) {
      this.quantity = quantity;
      return getThis();
   }

   public DropAttributesBuilder setMeso(Integer meso) {
      this.meso = meso;
      return getThis();
   }

   public DropAttributesBuilder setDropType(Integer dropType) {
      this.dropType = dropType;
      return getThis();
   }

   public DropAttributesBuilder setDropX(Integer dropX) {
      this.dropX = dropX;
      return getThis();
   }

   public DropAttributesBuilder setDropY(Integer dropY) {
      this.dropY = dropY;
      return getThis();
   }

   public DropAttributesBuilder setOwnerId(Integer ownerId) {
      this.ownerId = ownerId;
      return getThis();
   }

   public DropAttributesBuilder setOwnerPartyId(Integer ownerPartyId) {
      this.ownerPartyId = ownerPartyId;
      return getThis();
   }

   public DropAttributesBuilder setDropTime(Long dropTime) {
      this.dropTime = dropTime;
      return getThis();
   }

   public DropAttributesBuilder setDropperUniqueId(Integer dropperUniqueId) {
      this.dropperUniqueId = dropperUniqueId;
      return getThis();
   }

   public DropAttributesBuilder setDropperX(Integer dropperX) {
      this.dropperX = dropperX;
      return getThis();
   }

   public DropAttributesBuilder setDropperY(Integer dropperY) {
      this.dropperY = dropperY;
      return getThis();
   }

   public DropAttributesBuilder setPlayerDrop(Boolean playerDrop) {
      this.playerDrop = playerDrop;
      return getThis();
   }

   public DropAttributesBuilder setMod(Byte mod) {
      this.mod = mod;
      return getThis();
   }
}
