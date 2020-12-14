package com.atlas.drg.model;

public record Drop(int id, int itemId, int quantity, int type, int x, int y, int ownerId, Integer ownerPartyId, long dropTime,
                   int dropperId, int dropperX, int dropperY, boolean playerDrop, byte mod) {
}
