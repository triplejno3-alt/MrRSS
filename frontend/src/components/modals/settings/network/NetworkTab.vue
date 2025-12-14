<script setup lang="ts">
import type { SettingsData } from '@/types/settings';
import NetworkSettings from './NetworkSettings.vue';
import ProxySettings from './ProxySettings.vue';

defineProps<{
  settings: SettingsData;
}>();

const emit = defineEmits<{
  'update:settings': [settings: SettingsData];
}>();

// Handler for settings updates from child components
function handleUpdateSettings(updatedSettings: SettingsData) {
  emit('update:settings', updatedSettings);
}
</script>

<template>
  <div class="space-y-4 sm:space-y-6">
    <NetworkSettings />

    <ProxySettings :settings="settings" @update:settings="handleUpdateSettings" />
  </div>
</template>

<style scoped>
@reference "../../../../style.css";
</style>
