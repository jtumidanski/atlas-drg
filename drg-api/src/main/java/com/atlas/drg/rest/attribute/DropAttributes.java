package com.atlas.drg.rest.attribute;

import rest.AttributeResult;

public record DropAttributes(Integer worldId, Integer channelId, Integer mapId, Integer itemId, Integer quantity, Integer meso,
                             Integer dropType, Integer dropX, Integer dropY, Integer ownerId, Integer ownerPartyId,
                             Long dropTime, Integer dropperUniqueId, Integer dropperX, Integer dropperY, Boolean playerDrop,
                             Byte mod) implements AttributeResult {
}
