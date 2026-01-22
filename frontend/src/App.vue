<template>
  <div class="relative w-full h-full">
    <img
        class="absolute z-[1] left-0 top-0 w-full h-full object-cover"
        src="./assets/images/background.jpg"
        alt="aion2">
    <div class="absolute z-[2] left-0 top-0 w-full h-full bg-black/70"/>
    <div class="absolute z-[3] left-0 top-0 w-full h-full">
      <Home/>
    </div>
  </div>
</template>

<script lang="ts" setup>
import Home from "./components/Home.vue";
</script>
