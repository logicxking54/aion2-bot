<template>
  <div class="p-4 flex flex-col h-full">
    <div class="grow min-h-0 bg-black/50 mb-3 rounded-lg">
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
    <button
        v-if="!running"
        class="bg-green-700 py-3 px-10 rounded-lg"
        @click="onStart">
      เริ่มทำงาน
    </button>
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
import {IsRunning, Log, Start, Stop} from '../../wailsjs/go/main/App'

let running = ref(false)
let appLogs = ref<{
  date: string
  content: string
}[]>([])

const onStart = () => {
  Start()
}

const onStop = () => {
  Stop()
}

onMounted(() => {
  setInterval(async () => {
    appLogs.value = await Log() as any
    running.value = await IsRunning()

    console.log(appLogs)
  }, 500)
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