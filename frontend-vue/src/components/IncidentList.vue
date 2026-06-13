<script setup lang="ts">
import type { Incident } from '@/types';
import IncidentListItem from './IncidentListItem.vue';

const props = defineProps<{
    incidents : Array<Incident>
}>()
</script>

<template>
    <ul>
        <IncidentListItem
        v-for="inc in incidents"
        :key="inc.id"
        :inc="inc"
        />
    </ul>
</template>

<style>

</style>