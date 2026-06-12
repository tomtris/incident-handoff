<script setup lang="ts">
import type { Incident, TimelineEntry } from '@/types';
import { computed } from 'vue';

const props = defineProps<{inc : Incident}>()
// const entries = props.inc.entries
const actions = computed(()=>{
    return props.inc.entries.filter(n => n.type == "action")
})
const openQuestions = computed(() => {
    return props.inc.entries.filter(n => n.type == "open_question")
})
</script>

<template>
    <p v-if="actions.length">Actions taken:</p>
    <p v-for="action in actions" :key="action.id">{{ action.text }}</p>
    <p v-if="openQuestions.length">Last Open Question - Start Here</p>
    <p>{{ openQuestions[openQuestions.length - 1]?.text }}</p>
</template>
