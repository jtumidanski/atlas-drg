package com.atlas.drg.model;

public record MonsterDrop(int monsterId, int itemId, int maximumQuantity, int minimumQuantity, int chance) {
}
