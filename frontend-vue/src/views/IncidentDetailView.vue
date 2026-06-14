<script setup lang="ts">
import { addEntry, getIncident, logout, updateIncident, whoAmI } from '@/api';
import IncidentDetailBrief from '@/components/IncidentDetailBrief.vue';
import IncidentDetailAddEntry from '@/components/IncidentDetailAddEntry.vue';
import IncidentDetailEntryList from '@/components/IncidentDetailEntryList.vue';
import IncidentDetailHeader from '@/components/IncidentDetailHeader.vue';
import type { Incident, IncidentStatus, Severity, TimelineEntry, TimelineEntryType, UserContext } from '@/types';
import { onMounted, ref, watchEffect, type Ref } from 'vue';
import { useRoute } from 'vue-router';
import AppHeader from '@/components/AppHeader.vue';
import SeverityBadge from '@/components/SeverityBadge.vue';
import StatusBadge from '@/components/StatusBadge.vue';

const route = useRoute()
const inc = ref<Incident>()
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

// async function handleUpdateOnCall(payload: { key: string, value: string }) {
//     if (await runUpdate(payload) == true) {
//         inc.value.on_call = payload.value
//     }
// }
// async function handleUpdateSeverity(payload: { key: string, value: string }) {
//     await runUpdate(payload)
//     inc.value.severity = payload.value as Severity
// }
// async function handleUpdateStatus(payload: { key: string, value: string }) {
//     await runUpdate(payload)
//     inc.value.status = payload.value as IncidentStatus
// }

// async function runUpdate(payload: { key: string, value: string }) {
//     errIncidentLoadingMsg.value = ''
//     try {
//         await updateIncident(inc.value.id, { [payload.key]: payload.value })
//         return true
//     } catch (e) {
//         errIncidentLoadingMsg.value = (e as Error).message
//         return false
//     }
// }

// async function handleAddEntry(payload: {type:string, text:string}) {
//     errAddEntryMsg.value = ''
//     try {
//         const newEntry = await addEntry(inc.value.id, payload.type, payload.text)
//         inc.value.entries.push(newEntry)
//     } catch (e) {
//         errAddEntryMsg.value = (e as Error).message
//     }
// }
</script>

<template>
    <div>
        <AppHeader></AppHeader>
        <div class="page">
            <RouterLink :to="{name:'incidents'}" class="back mono">← Back to incident</RouterLink>
            
            <p v-if="errIncidentLoadingMsg" class="error"> {{ errIncidentLoadingMsg }}</p>
            <template v-else-if="inc">
                <div class="detail-head">
                    <div class="head-id-row">
                        <span class="detail-id mono">{{ inc.id }}</span>
                        <SeverityBadge :severity="inc.severity"></SeverityBadge>
                        <StatusBadge :status="inc.status"></StatusBadge>
                    </div>
                </div>
                <h1 class="detail-title"> {{ inc.title }}</h1>
                <div class="head-meta">
                    <span><span class="meta-key">service</span> <b class="mono">{{ inc.service }}</b></span>
                    <span><span class="meta-key">on-call</span> <b class="mono">{{ inc.on_call }}</b></span>
                    <span><span class="meta-key">opened by</span> <b class="mono">{{ inc.opened_by }}</b></span>
                    <!-- <span><span class="meta-key">elapsed</span> <b class="mono">{{ inc.elapsed }}</b></span> -->
                </div>
            </template>
        
            <div class="detail-grid">
                <div class="detail-main"></div>
            </div>
        </div>

    </div>
</template>