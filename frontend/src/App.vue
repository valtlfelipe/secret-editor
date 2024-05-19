<script lang="ts" setup>
import {
  SaveSecret,
  LoadSecret,
  GetPreferences,
  SetPreference,
  ExportSecret,
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
import { Button } from '@/components/ui/button';

const appVersion = APP_VERSION || '0';
const isInitializing = ref(true);
const isValid = ref(true);
const isDirty = ref(false);
const errorMessage = ref('');
const openNewSecret = ref(false);
const openSecretSelector = ref(false);
const openConfirm = ref(false);
let confirmIntendedAction: string = '';

const preferences: Ref<appType.Preferences> = ref({} as appType.Preferences);

const loadedSecret: Ref<appType.Secret> = ref({} as appType.Secret);

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
    confirmAction('new');
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
    confirmAction('open');
    return;
  }
  openSecretSelector.value = true;
});

EventsOn('command:export', () => {
  exportSecret();
});

async function exportSecret() {
  if (!loadedSecret.value.name) return;
  const result = await ExportSecret(loadedSecret.value.name, code.value);
  if (!result.success) {
    toast.error(result.error);
    return;
  }

  toast.success('Exported successfully.');
}

function confirmAction(action: string) {
  confirmIntendedAction = action;
  openConfirm.value = true;
}

function openSelectedSecret(item: appType.Secret) {
  loadSecret(item.arn);
  openSecretSelector.value = false;
}

async function submitNewSecret(item: appType.Secret) {
  loadSecret(item.arn);
  openNewSecret.value = false;
}

function confirmActionResult() {
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
        {
          from: 0,
          to: event.length,
          severity: 'error',
          message: e.message,
        },
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

  preferences.value = data;
  isInitializing.value = false;
}

let clickFooterTimes = 0;

function clickFooter() {
  clickFooterTimes++;
  if (clickFooterTimes === 5) {
    toast('You are a nice person! 😊');
    clickFooterTimes = 0;
  }
}

init();
</script>

<template>
  <div class="bg-gray-50">
    <Toaster />
    <NewSecret v-model:open="openNewSecret" @submit="submitNewSecret" />
    <SecretSelector
      v-model:open="openSecretSelector"
      @select="openSelectedSecret"
    />
    <WithoutSavingConfirm
      v-model:open="openConfirm"
      @confirm="confirmActionResult"
    />

    <div v-if="isInitializing">
      <Loading />
    </div>

    <div
      v-if="!isInitializing && loadedSecret.arn"
      class="flex flex-col h-screen"
    >
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
          <!-- <button class="text-white text-xs transition-all hover:font-bold">
            {{ preferences.provider.current }} [{{
              preferences.provider.awsProfile
            }}]
          </button> -->
          <p class="text-white text-xs">
            {{ isDirty ? `Not saved.` : 'Saved!' }}
          </p>
        </div>
      </footer>
    </div>

    <div
      v-if="!isInitializing && !loadedSecret.arn"
      class="flex flex-col h-screen"
    >
      <div
        class="mx-auto flex w-full h-screen flex-col justify-center space-y-6 sm:w-[350px]"
      >
        <div class="flex flex-col space-y-2 text-center">
          <h1 class="text-2xl font-semibold tracking-tight">Welcome!</h1>
          <Button
            @click="openNewSecret = true"
            variant="outline"
            class="text-muted-foreground"
          >
            Create a new Secret
            <kbd
              class="ml-4 pointer-events-none h-5 select-none gap-1 rounded border border-border bg-muted font-sans font-medium min-h-5 text-[14px] h-5 px-1 pointer-events-none h-5 select-none gap-1 rounded border bg-muted px-1.5 font-mono font-medium opacity-100"
              ><span class="text-xs">⌘</span>N
            </kbd>
          </Button>
          <Button
            @click="openSecretSelector = true"
            variant="outline"
            class="text-muted-foreground"
          >
            Open existing secret
            <kbd
              class="ml-4 pointer-events-none h-5 select-none gap-1 rounded border border-border bg-muted font-sans font-medium min-h-5 text-[14px] h-5 px-1 pointer-events-none h-5 select-none gap-1 rounded border bg-muted px-1.5 font-mono font-medium opacity-100"
              ><span class="text-xs">⌘</span>O
            </kbd>
          </Button>
        </div>
      </div>
      <footer class="mt-auto">
        <div class="py-3 text-center">
          <a
            href="#/"
            @click.prevent="clickFooter"
            class="text-muted-foreground text-xs"
            >Secret Editor v{{ appVersion }}</a
          >
        </div>
      </footer>
    </div>
  </div>
</template>

<style></style>
