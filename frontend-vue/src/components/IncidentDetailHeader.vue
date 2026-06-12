<script setup lang="ts">
import { whoAmI } from '@/api';
import type { Incident, Severity, IncidentStatus, UserContext } from '@/types';
import { onMount } from 'nanostores';
import { computed, onMounted, ref, watchEffect } from 'vue';

const props = defineProps<{
    inc: Incident,
    myIdentity: UserContext | undefined,
    currentOncall: string,
}>()
const emit = defineEmits<{
    updateSeverity: [value: Severity]
    updateStatus: [value: IncidentStatus]
    updateOnCall: [value: string]
}>()

const canEdit = computed(() => {
    return props.myIdentity?.role == "admin" || props.myIdentity?.username == props.currentOncall
})
</script>

<template>
    <h1>{{ inc.title }}</h1>

    <select :value="inc.severity" :disabled="!canEdit" @change="emit('updateSeverity', ($event.target as HTMLSelectElement).value as Severity)">
        <option>SEV1</option>
        <option>SEV2</option>
        <option>SEV3</option>
    </select>

    <select :value="inc.status" :disabled="!canEdit" @change="emit('updateStatus', ($event.target as HTMLSelectElement).value as IncidentStatus)">
        <option value="triggered">Triggered</option>
        <option value="acknowledged">Acknowledged</option>
        <option value="investigating">Investigating</option>
        <option value="mitigated">Mitigated</option>
        <option value="resolved">Resolved</option>
    </select>

    <input :value="inc.on_call" :disabled="!canEdit" @change="emit('updateOnCall', ($event.target as HTMLInputElement).value)" />

    <h3>{{ inc.service }}</h3>
    <p>{{ inc.id }}</p>
    <p>{{ inc.opened_by }}</p>
    <p>{{ inc.created_at }}</p>
</template>