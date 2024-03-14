<script setup lang="ts">
import {
  CommandDialog,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
} from '@/components/ui/command';
import { ref, watch } from 'vue';
import { toast } from 'vue-sonner';
import { GetSecrets } from '../../../../wailsjs/go/main/App';
import type { Ref } from 'vue';
import type { types as appType } from '../../../../wailsjs/go/models';

const secretsList: Ref<appType.Secret[]> = ref([]);
const open = defineModel<boolean>('open');
const emit = defineEmits(['select']);

async function getSecretsList() {
  const result = await GetSecrets();
  if (!result.success) {
    toast.error(result.error);
    return;
  }

  secretsList.value = result.result;
}

function openSelectedSecret(item: appType.Secret) {
  emit('select', item);
}

watch(open, (isOpen) => {
  if (isOpen) {
    getSecretsList();
  }
});
</script>

<template>
  <!-- tabindex="0" -->
  <CommandDialog v-model:open="open">
    <CommandInput
      placeholder="Open a secret 🤫 ..."
      @keydown.esc="open = false"
    />
    <CommandList>
      <CommandEmpty>No secret found.</CommandEmpty>
      <CommandGroup heading="Secrets">
        <CommandItem
          @select="openSelectedSecret(item)"
          v-for="item in secretsList"
          :value="item.arn"
        >
          {{ item.name }}
        </CommandItem>
      </CommandGroup>
    </CommandList>
  </CommandDialog>
</template>
