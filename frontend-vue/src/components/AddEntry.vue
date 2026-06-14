<script setup lang="ts">
import type { TimelineEntry, TimelineEntryType } from '@/types';
import { ref } from 'vue';

const newType = ref<TimelineEntryType>("observation")
const newText = ref("")

const props = defineProps<{
    errAddEntryMsg: string
}>()
const emit = defineEmits<{
    addEntry: [payload: {type: TimelineEntryType, text: string}]
}>()

function onClick() {
    emit('addEntry', {type: newType.value, text: newText.value})
}
</script>

<template>
    <div>
        <p class="error">{{ errAddEntryMsg }}</p>
        <h3 class="panel-title">Log an entry</h3>
        <div class="field">
            <label class="field-label">Type</label>
            <select class="select" v-model="newType">
            <option value="observation">Observation</option>
            <option value="action">Action</option>
            <option value="discovery">Discovery</option>
            <option value="open_question">Open question</option>
            <option value="state_change">State change</option>
            </select>
        </div>
        <div class="field">
            <label class="field-label">What happened</label>
            <textarea
            class="textarea"
            v-model="newText"
            placeholder="Describe the action, finding, or question"
            ></textarea>
        </div>
        <div class="row">
            <span class="spacer"></span>
            <button class="btn btn-primary" @click="onClick">Add to timeline</button>
        </div>
    </div>
</template>

<style scoped>
</style>