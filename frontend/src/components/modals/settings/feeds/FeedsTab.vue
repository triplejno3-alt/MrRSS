<script setup lang="ts">
import DataManagementSettings from './DataManagementSettings.vue';
import FeedManagementSettings from './FeedManagementSettings.vue';
import type { Feed } from '@/types/models';

const emit = defineEmits<{
  'import-opml': [];
  'export-opml': [];
  'cleanup-database': [];
  'add-feed': [];
  'edit-feed': [feed: Feed];
  'delete-feed': [id: number];
  'batch-delete': [ids: number[]];
  'batch-move': [ids: number[]];
  'discover-all': [];
}>();

// Event handlers that pass through to parent
function handleImportOPML() {
  emit('import-opml');
}

function handleExportOPML() {
  emit('export-opml');
}

function handleCleanupDatabase() {
  emit('cleanup-database');
}

function handleDiscoverAll() {
  emit('discover-all');
}

function handleAddFeed() {
  emit('add-feed');
}

function handleEditFeed(feed: Feed) {
  emit('edit-feed', feed);
}

function handleDeleteFeed(id: number) {
  emit('delete-feed', id);
}

function handleBatchDelete(ids: number[]) {
  emit('batch-delete', ids);
}

function handleBatchMove(ids: number[]) {
  emit('batch-move', ids);
}
</script>

<template>
  <div class="space-y-4 sm:space-y-6">
    <DataManagementSettings
      @import-opml="handleImportOPML"
      @export-opml="handleExportOPML"
      @cleanup-database="handleCleanupDatabase"
      @discover-all="handleDiscoverAll"
    />

    <FeedManagementSettings
      @add-feed="handleAddFeed"
      @edit-feed="handleEditFeed"
      @delete-feed="handleDeleteFeed"
      @batch-delete="handleBatchDelete"
      @batch-move="handleBatchMove"
    />
  </div>
</template>
