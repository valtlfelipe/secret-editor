<script lang="ts" setup>
import {
  SaveSecret,
  LoadSecret,
  GetPreferences,
  SetPreference,
} from '@wailsjs/go/main/App';
import { EventsOn, WindowSetTitle } from '@wailsjs/runtime';
import { Toaster } from '@/components/ui/sonner';
import { toast } from 'vue-sonner';

import { nextTick, ref, shallowRef } from 'vue';
import type { Ref } from 'vue';
import type { types as appType } from '@wailsjs/go/models';

import { Codemirror } from 'vue-codemirror';
import { json } from '@codemirror/lang-json';
import { oneDark } from '@codemirror/theme-one-dark';
import { setDiagnostics } from '@codemirror/lint';
import { NewSecret } from '@/components/app/new-secret';
import { SecretSelector } from '@/components/app/secret-selector';
import { WithoutSavingConfirm } from '@/components/app/without-saving-confirm';
import { Loading } from '@/components/ui/loading';

const isInitializing = ref(true);
const isValid = ref(true);
const isDirty = ref(false);
const errorMessage = ref('');
const openNewSecret = ref(false);
const openSecretSelector = ref(false);
const openConfirm = ref(false);
let confirmIntendedAction: string = '';

const loadedSecret: Ref<appType.Secret> = ref({
  arn: '',
  name: '',
  secret: '',
});

// Codemirror data
const code = ref(``);
const extensions = [json(), oneDark];
const view = shallowRef();
const handleReady = (payload: any) => {
  view.value = payload.view;
};

// commands events
EventsOn('command:new', async () => {
  if (openNewSecret.value) return;
  if (isDirty.value) {
    confirmIntendedAction = 'new';
    openConfirm.value = true;
    return;
  }
  openNewSecret.value = true;
});

EventsOn('command:save', async () => {
  saveLoadedSecret();
});

EventsOn('command:open', async () => {
  if (openSecretSelector.value) return;
  if (isDirty.value) {
    confirmIntendedAction = 'open';
    openConfirm.value = true;
    return;
  }
  openSecretSelector.value = true;
});

function openSelectedSecret(item: appType.Secret) {
  loadSecret(item.arn);
  openSecretSelector.value = false;
}

async function submitNewSecret(item: appType.Secret) {
  loadSecret(item.arn);
  openNewSecret.value = false;
}

function confirmAction() {
  if (confirmIntendedAction === '') return;

  switch (confirmIntendedAction) {
    case 'new':
      openNewSecret.value = true;
      break;

    case 'open':
      openSecretSelector.value = true;
      break;

    default:
      break;
  }

  confirmIntendedAction = '';
  openConfirm.value = false;
}

async function loadSecret(arn: string) {
  const result = await LoadSecret(arn);
  if (!result.success) {
    toast.error(result.error);
    return;
  }

  loadedSecret.value = result.result;
  code.value = JSON.stringify(JSON.parse(loadedSecret.value.secret), null, 2);
  nextTick(() => onCodeChange(code.value));

  WindowSetTitle(`Secret Editor :: ${loadedSecret.value.name}`);
  SetPreference('general.lastOpenedSecret', loadedSecret.value.name);
}

async function saveLoadedSecret() {
  if (!isValid.value || !loadedSecret.value.arn || !isDirty.value) return;
  const result = await SaveSecret(loadedSecret.value.arn, code.value);
  if (!result.success) {
    toast.error(result.error);
    return;
  }
  loadedSecret.value.secret = code.value;
  isDirty.value = false;
}

function onCodeChange(event: string) {
  isDirty.value = !(event === loadedSecret.value.secret);
  try {
    JSON.parse(event);
  } catch (e: any) {
    view.value.dispatch(
      setDiagnostics(view.value.state, [
        { from: 0, to: event.length, severity: 'error', message: e.message },
      ])
    );
    errorMessage.value = e.message;
    isValid.value = false;
    return;
  }

  isValid.value = true;
  view.value.dispatch(setDiagnostics(view.value.state, []));
}

async function init() {
  const data = await GetPreferences();

  if (data.general.lastOpenedSecret) {
    await loadSecret(data.general.lastOpenedSecret);
  }

  isInitializing.value = false;
}

init();
</script>

<template>
  <div>
    <Toaster />
    <NewSecret v-model:open="openNewSecret" @submit="submitNewSecret" />
    <SecretSelector
      v-model:open="openSecretSelector"
      @select="openSelectedSecret"
    />
    <WithoutSavingConfirm v-model:open="openConfirm" @confirm="confirmAction" />

    <div v-if="isInitializing">
      <Loading />
    </div>

    <div v-else class="flex flex-col h-screen">
      <header class="flex items-center justify-between space-y-2 p-4">
        <div>
          <h2 class="text-2xl font-bold tracking-tight">
            {{ loadedSecret.name }}
          </h2>
          <p class="text-muted-foreground">
            {{ loadedSecret.arn }}
          </p>
        </div>
        <div class="flex items-center space-x-2">
          <!-- <Button>Click me</Button> -->
        </div>
      </header>
      <main class="flex-grow" style="overflow: auto">
        <codemirror
          v-model="code"
          placeholder="Secret JSON goes here..."
          :style="{
            height: '100%',
            width: '100%',
            outline: 'none',
            overflow: 'auto',
          }"
          :autofocus="true"
          :indent-with-tab="true"
          :tab-size="2"
          :extensions="extensions"
          @ready="handleReady"
          @change="onCodeChange"
        />
      </main>
      <footer>
        <div
          class="w-full mx-auto p-2 px-4 md:flex md:items-center md:justify-between bg-primary"
        >
          <p class="text-white text-xs">
            {{ isValid ? '✅ Valid' : `⚠️ ${errorMessage}` }}
          </p>
          <p class="text-white text-xs">
            {{ isDirty ? `Not saved.` : 'Saved!' }}
          </p>
        </div>
      </footer>
    </div>
  </div>
</template>

<style></style>
