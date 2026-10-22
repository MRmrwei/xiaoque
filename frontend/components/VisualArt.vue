<template>
  <div class="visual-art" :class="`visual-art--${variant}`">
    <!-- Orbit -->
    <svg v-if="variant === 'orbit'" viewBox="0 0 600 450" preserveAspectRatio="xMidYMid slice">
      <rect width="600" height="450" :fill="c.bg" />
      <circle cx="350" cy="225" r="205" fill="none" :stroke="c.paper" stroke-opacity=".22" />
      <circle cx="350" cy="225" r="152" fill="none" :stroke="c.paper" stroke-opacity=".34" />
      <circle cx="350" cy="225" r="98" fill="none" :stroke="c.paper" stroke-opacity=".5" />
      <circle cx="350" cy="225" r="46" :fill="c.paper" />
      <circle cx="350" cy="225" r="16" :fill="c.ink" />
      <circle cx="456" cy="142" r="13" :fill="c.accent" :stroke="c.bg" stroke-width="4" />
      <circle cx="238" cy="328" r="7" :fill="c.paper" />
      <path d="M42 390h72" :stroke="c.paper" stroke-opacity=".45" />
      <path d="M42 402h40" :stroke="c.paper" stroke-opacity=".25" />
      <g :fill="c.paper" opacity=".55">
        <circle cx="86" cy="92" r="3" />
        <circle cx="500" cy="350" r="3" />
        <circle cx="530" cy="92" r="3" />
      </g>
      <path d="M0 0l600 450" :stroke="c.paper" stroke-opacity=".08" />
    </svg>

    <!-- Arch -->
    <svg v-else-if="variant === 'arch'" viewBox="0 0 600 450" preserveAspectRatio="xMidYMid slice">
      <rect width="600" height="450" :fill="c.bg" />
      <circle cx="474" cy="112" r="58" :fill="c.accent" />
      <path d="M92 450V268a136 136 0 0 1 272 0v182Z" :fill="c.paper" fill-opacity=".14" />
      <path d="M150 450V286a88 88 0 0 1 176 0v164Z" :fill="c.paper" fill-opacity=".2" />
      <path d="M208 450V304a44 44 0 0 1 88 0v146Z" :fill="c.accent" />
      <path d="M420 450V352a48 48 0 0 1 96 0v98Z" :fill="c.paper" fill-opacity=".5" />
      <path d="M0 450h600" :stroke="c.paper" stroke-opacity=".3" />
      <g :stroke="c.paper" stroke-opacity=".2">
        <path d="M40 80h110" />
        <path d="M40 94h64" />
      </g>
      <circle cx="304" cy="170" r="4" :fill="c.paper" />
    </svg>

    <!-- Grid -->
    <svg v-else-if="variant === 'grid'" viewBox="0 0 600 450" preserveAspectRatio="xMidYMid slice">
      <defs>
        <pattern id="va-grid" width="30" height="30" patternUnits="userSpaceOnUse">
          <path d="M30 0H0v30" fill="none" :stroke="c.ink" stroke-opacity=".12" />
        </pattern>
      </defs>
      <rect width="600" height="450" :fill="c.bg" />
      <rect width="600" height="450" fill="url(#va-grid)" />
      <circle cx="408" cy="186" r="128" :fill="c.accent" />
      <rect x="92" y="284" width="146" height="146" rx="20" :fill="c.ink" />
      <circle cx="165" cy="357" r="34" :fill="c.bg" />
      <path d="M92 230h146" :stroke="c.ink" stroke-width="2" />
      <path d="M392 350h116" :stroke="c.ink" stroke-width="2" />
      <g :fill="c.ink" opacity=".5">
        <circle cx="92" cy="94" r="3" />
        <circle cx="522" cy="94" r="3" />
        <circle cx="522" cy="408" r="3" />
      </g>
    </svg>

    <!-- Dunes -->
    <svg v-else-if="variant === 'dunes'" viewBox="0 0 600 450" preserveAspectRatio="xMidYMid slice">
      <defs>
        <linearGradient id="va-dune-a" x1="0" y1="0" x2="1" y2="1">
          <stop offset="0" :stop-color="c.accent" />
          <stop offset="1" :stop-color="c.accentDeep" />
        </linearGradient>
      </defs>
      <rect width="600" height="450" :fill="c.bg" />
      <circle cx="438" cy="120" r="62" :fill="c.paper" fill-opacity=".85" />
      <path d="M0 322c94-38 168-40 246-8 92 38 168 42 354-30v166H0Z" fill="url(#va-dune-a)" />
      <path d="M0 362c112-34 198-20 268 10 94 40 190 18 332-44v122H0Z" :fill="c.ink" fill-opacity=".88" />
      <path d="M0 408c132-26 228-12 322 14 76 18 164 6 278-32v60H0Z" :fill="c.paper" fill-opacity=".55" />
      <g :stroke="c.paper" stroke-opacity=".55">
        <path d="M52 76h104" />
        <path d="M52 90h58" />
      </g>
      <path d="M466 234c26-18 46-20 78-8" :stroke="c.ink" stroke-opacity=".4" fill="none" />
    </svg>

    <!-- Waves -->
    <svg v-else-if="variant === 'waves'" viewBox="0 0 600 450" preserveAspectRatio="xMidYMid slice">
      <rect width="600" height="450" :fill="c.bg" />
      <circle cx="446" cy="128" r="54" :fill="c.accent" />
      <g fill="none" :stroke="c.paper" stroke-opacity=".24" stroke-width="1.4">
        <path d="M-20 120c84-46 148-46 232 0s148 46 232 0 148-46 232 0" />
        <path d="M-20 160c84-46 148-46 232 0s148 46 232 0 148-46 232 0" />
        <path d="M-20 200c84-46 148-46 232 0s148 46 232 0 148-46 232 0" />
        <path d="M-20 240c84-46 148-46 232 0s148 46 232 0 148-46 232 0" />
        <path d="M-20 280c84-46 148-46 232 0s148 46 232 0 148-46 232 0" />
        <path d="M-20 320c84-46 148-46 232 0s148 46 232 0 148-46 232 0" />
        <path d="M-20 360c84-46 148-46 232 0s148 46 232 0 148-46 232 0" />
      </g>
      <path
        d="M-20 280c84-46 148-46 232 0s148 46 232 0 148-46 232 0"
        fill="none"
        :stroke="c.accent"
        stroke-width="4"
      />
      <circle cx="388" cy="302" r="6" :fill="c.paper" />
    </svg>

    <!-- Type -->
    <svg v-else-if="variant === 'type'" viewBox="0 0 600 450" preserveAspectRatio="xMidYMid slice">
      <rect width="600" height="450" :fill="c.bg" />
      <text
        x="300"
        y="336"
        text-anchor="middle"
        :fill="c.paper"
        font-family="'Instrument Serif','Noto Serif SC','Songti SC',serif"
        font-size="300"
        opacity=".16"
      >
        鹊
      </text>
      <text
        x="294"
        y="326"
        text-anchor="middle"
        :fill="c.paper"
        font-family="'Instrument Serif','Noto Serif SC','Songti SC',serif"
        font-size="256"
      >
        鹊
      </text>
      <path d="M70 88h92" :stroke="c.paper" stroke-opacity=".55" />
      <text
        x="70"
        y="72"
        :fill="c.paper"
        font-family="'SFMono-Regular',monospace"
        font-size="12"
        letter-spacing="3"
        opacity=".65"
      >
        {{ brandLabel }}
      </text>
      <circle cx="70" cy="392" r="4" :fill="c.accent" />
    </svg>

    <!-- Stack -->
    <svg v-else-if="variant === 'stack'" viewBox="0 0 600 450" preserveAspectRatio="xMidYMid slice">
      <rect width="600" height="450" :fill="c.bg" />
      <g transform="rotate(-9 300 240)">
        <rect x="150" y="210" width="300" height="170" rx="18" :fill="c.ink" fill-opacity=".16" />
        <rect x="162" y="190" width="300" height="170" rx="18" :fill="c.ink" fill-opacity=".35" />
        <rect x="174" y="170" width="300" height="170" rx="18" :fill="c.paper" />
        <rect x="202" y="198" width="118" height="12" rx="6" :fill="c.ink" />
        <rect x="202" y="222" width="196" height="7" rx="3.5" :fill="c.ink" fill-opacity=".25" />
        <rect x="202" y="238" width="160" height="7" rx="3.5" :fill="c.ink" fill-opacity=".18" />
        <rect x="202" y="276" width="70" height="34" rx="10" :fill="c.accent" />
        <circle cx="420" cy="284" r="18" :fill="c.bg" :stroke="c.ink" stroke-opacity=".2" />
      </g>
      <circle cx="120" cy="96" r="26" :fill="c.accent" />
      <path d="M462 392h82" :stroke="c.paper" stroke-opacity=".45" />
      <path d="M462 404h44" :stroke="c.paper" stroke-opacity=".25" />
    </svg>

    <!-- Bloom -->
    <svg v-else viewBox="0 0 600 450" preserveAspectRatio="xMidYMid slice">
      <rect width="600" height="450" :fill="c.bg" />
      <g transform="translate(300 225)">
        <ellipse
          v-for="n in 14"
          :key="n"
          :cx="Math.cos(((n - 1) / 14) * Math.PI * 2) * 92"
          :cy="Math.sin(((n - 1) / 14) * Math.PI * 2) * 92"
          rx="78"
          ry="26"
          :fill="c.accent"
          :opacity="0.12 + (n % 4) * 0.075"
          :transform="`rotate(${((n - 1) / 14) * 360})`"
        />
        <circle r="52" :fill="c.paper" />
        <circle r="18" :fill="c.ink" />
      </g>
      <circle cx="92" cy="88" r="5" :fill="c.paper" />
      <circle cx="520" cy="366" r="5" :fill="c.accent" />
      <path d="M58 380h88" :stroke="c.paper" stroke-opacity=".4" />
    </svg>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useSiteChrome } from '~/composables/useSiteChrome'

const props = withDefaults(
  defineProps<{
    variant?: string
    tone?: string
  }>(),
  { variant: 'orbit', tone: '' }
)

const { settings } = useSiteChrome()
const brandLabel = computed(() => {
  const name = settings.value.siteName.en || settings.value.siteName.zh || 'Magpie'
  return `${name.toUpperCase()} / ${new Date().getFullYear()}`
})

interface Tone {
  bg: string
  accent: string
  accentDeep: string
  ink: string
  paper: string
}

const tones: Record<string, Tone> = {
  accent: { bg: '#E8503A', accent: '#141518', accentDeep: '#A72E1D', ink: '#141518', paper: '#F6F2EA' },
  ink: { bg: '#141518', accent: '#E8503A', accentDeep: '#C93B26', ink: '#F6F2EA', paper: '#F6F2EA' },
  sage: { bg: '#3F574C', accent: '#E8503A', accentDeep: '#C93B26', ink: '#141518', paper: '#F4EFE5' },
  cream: { bg: '#EFE6D6', accent: '#E8503A', accentDeep: '#C93B26', ink: '#141518', paper: '#FBF9F5' },
  sand: { bg: '#E6D8C3', accent: '#E8503A', accentDeep: '#C93B26', ink: '#141518', paper: '#F6F2EA' }
}

const defaultTone: Record<string, string> = {
  orbit: 'accent',
  arch: 'ink',
  grid: 'cream',
  dunes: 'sand',
  waves: 'ink',
  type: 'accent',
  stack: 'sage',
  bloom: 'cream'
}

const c = computed<Tone>(
  () => tones[props.tone || defaultTone[props.variant] || 'accent'] ?? tones.accent!
)
</script>

<style scoped>
.visual-art {
  width: 100%;
  height: 100%;
  overflow: hidden;
  background: var(--paper-deep);
}

.visual-art svg {
  width: 100%;
  height: 100%;
  display: block;
}
</style>
