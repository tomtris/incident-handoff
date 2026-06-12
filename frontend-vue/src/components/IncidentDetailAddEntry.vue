<script setup lang="ts">
import type { Incident } from '@/types';
import { ref } from 'vue'

const props = defineProps<{
    inc: Incident,
    err: string,
}>()

const entryType = ref('')
const entryText = ref('')
const emit = defineEmits<{
    submit: [payload: {incidentID: string, type: string, text: string}]
}>()

async function onSubmit() {
    emit("submit", {incidentID: props.inc.id, type: entryType.value, text: entryText.value})
}
</script>

<template>
    <p v-if="err" class="error">{{ err }}</p>
    <form @submit.prevent="onSubmit">
        <h1>Action Timeline</h1>
        <label>Type
            <input v-model="entryType">
        </label>
        <label>Text
            <input v-model="entryText">
        </label>
        <button type="submit">Log Entry</button>
    </form>
</template>
