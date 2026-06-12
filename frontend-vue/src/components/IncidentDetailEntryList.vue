<script setup lang="ts">
import type { Incident } from '@/types';

const props = defineProps<{inc : Incident}>()
</script>

<template>
    <div v-for="entry in inc.entries">
        <p>{{ entry.created_at }} {{ entry.author }} {{ entry.type }}</p>
        <p>{{ entry.text }}</p>
    </div>
</template>
