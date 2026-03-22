<template>
  <div class="p-4 flex flex-col h-full">
    <div class="grid grid-cols-2 gap-x-2 mb-3">
      <div
          class="bg-black/50 text-left p-4 rounded-lg">
        <p
            class="text-xs text-white"
        >
          เวลา: {{ timeStr || 'รอเริ่มงาน' }}
        </p>
      </div>

      <div class="grid grid-cols-3 gap-x-2">
        <div
            class="bg-black/50 text-left p-4 rounded-lg">
          <p
              class="text-xs text-white"
          >
            E: {{ eCooldown }}
          </p>
        </div>
        <div
            class="bg-black/50 text-left p-4 rounded-lg">
          <p
              class="text-xs text-white"
          >
            R: {{ rCooldown }}
          </p>
        </div>
        <div
            class="bg-black/50 text-left p-4 rounded-lg">
          <p
              class="text-xs text-white"
          >
            L: {{ lCooldown }}
          </p>
        </div>
      </div>
    </div>
    <div class="grow min-h-0 bg-black/50 mb-4 rounded-lg">
      <div
          id="log-container"
          class="p-4 max-h-full w-full overflow-auto text-left"
      >
        <p
            v-for="(l, i) in appLogs"
            :key="i"
            class="mb-1 text-xs text-white"
        >
          {{ l.date }}: {{ l.content }}
        </p>
      </div>
    </div>
    <div
        v-if="!running"
        class="flex items-center justify-between space-x-2"
    >
      <button
          class="bg-green-700 py-3 w-full px-10 rounded-lg"
          @click="onStartBot">
        ปล่อยบอท
      </button>
      <button
          class="bg-green-700 py-3 w-full px-10 rounded-lg"
          @click="onStartAutoKey">
        รัน auto key
      </button>
    </div>
    <button
        v-else-if="running"
        class="bg-red-700 py-3 px-10 rounded-lg"
        @click="onStop">
      หยุดการทำงาน
    </button>
  </div>
</template>

<script setup lang="ts">
import {nextTick, onMounted, ref, watch} from 'vue'
import {
  ECooldown,
  IsRunning,
  LCooldown,
  Log,
  RCooldown,
  StartAutoKey,
  StartBot,
  Stop,
  TimeStr
} from '../../wailsjs/go/main/App'

let timeStr = ref('')
let eCooldown = ref('')
let rCooldown = ref('')
let lCooldown = ref('')
let running = ref(false)
let appLogs = ref<{
  date: string
  content: string
}[]>([])

const onStartBot = () => {
  StartBot()
}

const onStartAutoKey = () => {
  StartAutoKey()
}

const onStop = () => {
  Stop()
}

onMounted(() => {
  setInterval(async () => {
    timeStr.value = await TimeStr()
    eCooldown.value = await ECooldown()
    lCooldown.value = await RCooldown()
    rCooldown.value = await LCooldown()
    appLogs.value = await Log() as any
    running.value = await IsRunning()

    console.log(appLogs)
  }, 50)
})

watch(() => appLogs.value.length, async (v, o) => {
  if (v !== o) {
    await nextTick()

    const c = document.getElementById('log-container')

    c?.scroll({
      top: c?.scrollHeight,
      behavior: 'smooth',
    })
  }
})
</script>