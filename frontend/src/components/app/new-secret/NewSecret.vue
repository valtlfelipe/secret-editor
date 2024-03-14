<script setup lang="ts">
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import { CreateSecret } from '../../../../wailsjs/go/main/App';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Button } from '@/components/ui/button';
import { ref, watch } from 'vue';
import { toast } from 'vue-sonner';

const newSecretName = ref('');
const open = defineModel<boolean>('open');
const emit = defineEmits(['submit']);

async function submitNewSecret() {
  const result = await CreateSecret(newSecretName.value, '{}');
  if (!result.success) {
    toast.error(result.error);
    return;
  }

  emit('submit', result.result);

  newSecretName.value = '';
  toast(`'${result.result.name}' created successfully.`);
}

watch(open, (isOpen) => {
  if (!isOpen) {
    newSecretName.value = '';
  }
});
</script>

<template>
  <Dialog v-model:open="open">
    <DialogContent class="sm:max-w-[450px]">
      <DialogHeader>
        <DialogTitle> New Secret </DialogTitle>
        <DialogDescription> Create a new empty secret. </DialogDescription>
      </DialogHeader>
      <form @submit.prevent="submitNewSecret">
        <div class="grid gap-4 py-4">
          <div class="grid grid-cols-4 items-center gap-4">
            <Label for="name" class="text-right"> Name </Label>
            <Input id="name" v-model="newSecretName" class="col-span-3" />
          </div>
        </div>
        <DialogFooter>
          <Button type="submit"> Create </Button>
        </DialogFooter>
      </form>
    </DialogContent>
  </Dialog>
</template>
