<script setup lang="ts">
import { addEntry, getIncident, updateIncident, whoAmI } from '@/api';
import IncidentDetailBrief from '@/components/IncidentDetailBrief.vue';
import IncidentDetailAddEntry from '@/components/IncidentDetailAddEntry.vue';
import IncidentDetailEntryList from '@/components/IncidentDetailEntryList.vue';
import IncidentDetailHeader from '@/components/IncidentDetailHeader.vue';
import type { Incident, IncidentStatus, Severity, TimelineEntry, TimelineEntryType, UserContext } from '@/types';
import { onMounted, ref, watchEffect, type Ref } from 'vue';
import { useRoute } from 'vue-router';

const route = useRoute()
const inc = ref<Incident>() as Ref<Incident>
const errIncidentLoadingMsg = ref('')
const errWhoAmI = ref('')
const errAddEntryMsg = ref('')
const myIdentity = ref<UserContext>({ id: "", username: "", role: "" })

onMounted(async() => {
    try {
        myIdentity.value = await whoAmI()
    } catch (e) {
        errWhoAmI.value = (e as Error).message
    }
    try {
        inc.value = await getIncident(route.params.id as string)
    } catch (e) {
        errIncidentLoadingMsg.value = (e as Error).message
    }
})

async function handleUpdateOnCall(payload: { key: string, value: string }) {
    if (await runUpdate(payload) == true) {
        inc.value.on_call = payload.value
    }
}
async function handleUpdateSeverity(payload: { key: string, value: string }) {
    await runUpdate(payload)
    inc.value.severity = payload.value as Severity
}
async function handleUpdateStatus(payload: { key: string, value: string }) {
    await runUpdate(payload)
    inc.value.status = payload.value as IncidentStatus
}

async function runUpdate(payload: { key: string, value: string }) {
    errIncidentLoadingMsg.value = ''
    try {
        await updateIncident(inc.value.id, { [payload.key]: payload.value })
        return true
    } catch (e) {
        errIncidentLoadingMsg.value = (e as Error).message
        return false
    }
}

async function handleAddEntry(payload: {type:string, text:string}) {
    errAddEntryMsg.value = ''
    try {
        const newEntry = await addEntry(inc.value.id, payload.type, payload.text)
        inc.value.entries.push(newEntry)
    } catch (e) {
        errAddEntryMsg.value = (e as Error).message
    }
}
</script>

<template>
    <RouterLink :to="{name: 'entry'}">Back</RouterLink>
    
    <div v-if="errIncidentLoadingMsg">
        <p>Failed to load incident:</p>
        <p>{{ errIncidentLoadingMsg }}</p>
        <p>Please try again</p>
    </div>
    <div v-else-if="errWhoAmI">
        <p>Failed to get your Identity:</p>
        <p>{{ errWhoAmI }}</p>
        <p>Please try again</p>
    </div>

    <div v-else-if="inc">
        <div>
            <IncidentDetailHeader
                :key="inc.id"
                :inc="inc"
                :myIdentity="myIdentity"
                @update-on-call="handleUpdateOnCall"
                @update-severity="handleUpdateSeverity"
                @update-status="handleUpdateStatus"
            />
        </div>
        
        <div>
            <IncidentDetailBrief
                :key="inc.id"
                :inc="inc"
            />
        </div>
        
        <div>
            <IncidentDetailAddEntry
                :key="inc.id"
                :err="errAddEntryMsg" 
                @submit="handleAddEntry"
            />
        </div>
        
        <div>
            <IncidentDetailEntryList
                :key="inc.id"
                :inc="inc"
            />
        </div>
    </div>
</template>