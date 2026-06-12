<script setup lang="ts">
import { addEntry, getIncident, whoAmI } from '@/api';
import IncidentDetailBrief from '@/components/IncidentDetailBrief.vue';
import IncidentDetailAddEntry from '@/components/IncidentDetailAddEntry.vue';
import IncidentDetailEntryList from '@/components/IncidentDetailEntryList.vue';
import IncidentDetailHeader from '@/components/IncidentDetailHeader.vue';
import type { Incident, TimelineEntry, TimelineEntryType, UserContext } from '@/types';
import { onMounted, ref, watchEffect } from 'vue';
import { useRoute } from 'vue-router';

const route = useRoute()
const inc = ref<Incident | null>(null)
const errIncidentLoadingMsg = ref('')
const errWhoAmI = ref('')
const errAddEntryMsg = ref('')
const myIdentity = ref<UserContext>()

onMounted(async() => {
    try {
        myIdentity.value = await whoAmI()
    } catch (e) {
        errWhoAmI.value = (e as Error).message
    }
})

watchEffect(async () => {
    try {
        inc.value = await getIncident(route.params.id as string)
    } catch (e) {
        errIncidentLoadingMsg.value = (e as Error).message
    }
})

async function handleAddEntry(payload: {incidentID: string, type:string, text:string}) {
    try {
        const newEntry = await addEntry(payload.incidentID, payload.type, payload.text)
        inc.value?.entries.push(newEntry)
    } catch (e) {
        errAddEntryMsg.value = (e as Error).message
    }
}
</script>

<template>
    <RouterLink :to="{name: 'incidents'}">Back</RouterLink>
    
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
                :currentOncall="inc.on_call"
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
                :inc="inc"
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