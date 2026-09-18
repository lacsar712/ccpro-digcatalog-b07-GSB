<template>
  <div>
    <div class="toolbar">
      <div>
        <h2 class="page-title">工作台</h2>
        <p class="page-sub">
          {{ data.today || '今日' }}（{{ data.timezone || 'UTC+8' }}）· 今日新增与最近登记动态
        </p>
      </div>
      <button class="btn secondary" @click="load">刷新</button>
    </div>

    <div v-if="error" class="error card">{{ error }}</div>

    <div class="stats">
      <router-link to="/finds" class="stat card highlight">
        <div class="label">今日新增文物</div>
        <div class="value">{{ loaded ? data.todayFindCount : '-' }}</div>
        <div class="hint">前往出土文物 →</div>
      </router-link>
      <router-link to="/sites" class="stat card">
        <div class="label">发掘工地</div>
        <div class="value">{{ loaded ? data.siteCount : '-' }}</div>
        <div class="hint">前往发掘工地 →</div>
      </router-link>
      <router-link to="/units" class="stat card">
        <div class="label">探方/发掘单位</div>
        <div class="value">{{ loaded ? data.unitCount : '-' }}</div>
        <div class="hint">前往探方单位 →</div>
      </router-link>
      <router-link to="/finds" class="stat card">
        <div class="label">出土文物总数</div>
        <div class="value">{{ loaded ? data.findCount : '-' }}</div>
        <div class="hint">前往出土文物 →</div>
      </router-link>
    </div>

    <div class="card">
      <div class="recent-head">
        <h3>最近登记文物</h3>
        <router-link to="/finds" class="more">查看全部 →</router-link>
      </div>
      <table class="table" v-if="data.recentFinds?.length">
        <thead>
          <tr>
            <th>登记号</th>
            <th>探方</th>
            <th>器物类型</th>
            <th>材质</th>
            <th>完整度</th>
            <th>登记时间</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="item in data.recentFinds" :key="item.id" @click="goFinds" class="clickable">
            <td>{{ item.registerNo }}</td>
            <td>{{ item.unit?.site?.name || '' }} {{ item.unit?.code || '-' }}</td>
            <td><span class="tag">{{ item.artifactType }}</span></td>
            <td>{{ item.materialName || item.material?.name || '-' }}</td>
            <td>{{ item.completeness || '-' }}</td>
            <td>{{ formatTime(item.createdAt) }}</td>
          </tr>
        </tbody>
      </table>
      <p v-else class="page-sub">今日暂无登记记录，去<router-link to="/finds">出土文物</router-link>新增一条吧</p>
    </div>
  </div>
</template>

<script setup>
import { onMounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import api from '../api/http'

const router = useRouter()

const data = reactive({
  timezone: '',
  today: '',
  todayFindCount: 0,
  recentFinds: [],
  siteCount: 0,
  unitCount: 0,
  findCount: 0
})
const loaded = ref(false)
const error = ref('')

function formatTime(v) {
  if (!v) return '-'
  const d = new Date(v)
  if (Number.isNaN(d.getTime())) return '-'
  // 与后端统计口径一致，按东八区展示
  return d.toLocaleString('zh-CN', { timeZone: 'Asia/Shanghai', hour12: false })
}

function goFinds() {
  router.push('/finds')
}

async function load() {
  error.value = ''
  try {
    const { data: res } = await api.get('/workspace/today')
    Object.assign(data, res)
    loaded.value = true
  } catch (e) {
    error.value = e.response?.data?.error || '加载失败'
  }
}

onMounted(load)
</script>

<style scoped>
.stats {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 1rem;
  margin-bottom: 1rem;
}

.stat {
  display: block;
  transition: transform 0.12s ease, box-shadow 0.12s ease;
}

.stat:hover {
  transform: translateY(-2px);
  box-shadow: 0 12px 28px rgba(60, 40, 20, 0.14);
}

.stat.highlight {
  border-color: var(--accent);
}

.stat .label {
  color: var(--muted);
  font-size: 0.9rem;
}

.stat .value {
  font-size: 2rem;
  font-weight: 700;
  margin-top: 0.35rem;
  color: var(--accent);
}

.stat .hint {
  margin-top: 0.35rem;
  font-size: 0.8rem;
  color: var(--muted);
}

.recent-head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0.75rem;
}

.recent-head h3 {
  margin: 0;
}

.more {
  font-size: 0.88rem;
  color: var(--accent);
}

.clickable {
  cursor: pointer;
}

.clickable:hover td {
  background: rgba(196, 165, 116, 0.12);
}

@media (max-width: 1000px) {
  .stats {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 600px) {
  .stats {
    grid-template-columns: 1fr;
  }
}
</style>
