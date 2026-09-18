<template>
  <div>
    <div class="toolbar">
      <div>
        <h2 class="page-title">工作台</h2>
        <p class="page-sub">今日工作概览（按东八区自然日统计）</p>
      </div>
      <button class="btn secondary" @click="load">刷新</button>
    </div>

    <div v-if="error" class="error card">{{ error }}</div>

    <div class="stats">
      <div class="stat card clickable" @click="go('finds')">
        <div class="label">今日新增文物</div>
        <div class="value">{{ data.todayFindCount ?? '-' }}</div>
        <div class="hint">进入出土文物 →</div>
      </div>
      <div class="stat card clickable" @click="go('sites')">
        <div class="label">发掘工地</div>
        <div class="value">{{ data.siteCount ?? '-' }}</div>
        <div class="hint">进入发掘工地 →</div>
      </div>
      <div class="stat card clickable" @click="go('units')">
        <div class="label">探方/发掘单位</div>
        <div class="value">{{ data.unitCount ?? '-' }}</div>
        <div class="hint">进入探方单位 →</div>
      </div>
      <div class="stat card clickable" @click="go('finds')">
        <div class="label">出土文物总数</div>
        <div class="value">{{ data.findCount ?? '-' }}</div>
        <div class="hint">进入出土文物 →</div>
      </div>
    </div>

    <div class="card">
      <div class="recent-head">
        <h3>最近登记文物</h3>
        <router-link class="more" :to="{ name: 'finds' }">查看全部 →</router-link>
      </div>
      <table class="table" v-if="data.recentFinds?.length">
        <thead>
          <tr>
            <th>登记号</th>
            <th>工地 / 探方</th>
            <th>器物类型</th>
            <th>材质</th>
            <th>完整度</th>
            <th>出土日期</th>
            <th>登记时间</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="item in data.recentFinds" :key="item.id">
            <td>{{ item.registerNo }}</td>
            <td>{{ item.siteName || '-' }} / {{ item.unitCode || '-' }}</td>
            <td><span class="tag">{{ item.artifactType }}</span></td>
            <td>{{ item.materialName || '-' }}</td>
            <td>{{ item.completeness || '-' }}</td>
            <td>{{ formatDate(item.findDate) }}</td>
            <td>{{ formatTime(item.createdAt) }}</td>
          </tr>
        </tbody>
      </table>
      <p v-else class="page-sub">
        暂无登记数据，去<router-link :to="{ name: 'finds' }">出土文物</router-link>新增一条吧
      </p>
    </div>
  </div>
</template>

<script setup>
import { onMounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import api from '../api/http'

const router = useRouter()

const data = reactive({
  todayFindCount: 0,
  siteCount: 0,
  unitCount: 0,
  findCount: 0,
  recentFinds: []
})
const error = ref('')

function go(name) {
  router.push({ name })
}

function formatDate(v) {
  if (!v) return '-'
  return String(v).slice(0, 10)
}

function formatTime(v) {
  if (!v) return '-'
  return String(v).slice(0, 16).replace('T', ' ')
}

async function load() {
  error.value = ''
  try {
    const { data: res } = await api.get('/workspace/today')
    Object.assign(data, res)
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
  margin-top: 0.4rem;
  font-size: 0.8rem;
  color: var(--muted);
}

.clickable {
  cursor: pointer;
  transition: transform 0.12s ease, box-shadow 0.12s ease;
}

.clickable:hover {
  transform: translateY(-2px);
  box-shadow: 0 12px 28px rgba(60, 40, 20, 0.14);
}

.clickable:hover .hint {
  color: var(--accent);
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
