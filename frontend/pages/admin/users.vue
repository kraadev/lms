<script setup lang="ts">
import { Users, Plus, Search, Pencil, Trash2, Shield, GraduationCap, School } from 'lucide-vue-next'
import { adminService } from '~/services/admin'
import type { User, UserRole } from '~/types'
import { formatDate } from '~/utils/formatters'

definePageMeta({
  middleware: ['auth', 'admin']
})

useSeoMeta({ title: 'Manajemen Pengguna - Admin' })

const toast = useToast()
const users = ref<User[]>([])
const isLoading = ref(true)
const error = ref<string | null>(null)

// Filters
const searchQuery = ref('')
const selectedRole = ref<UserRole | ''>('')

// Modals
const showModal = ref(false)
const editingUser = ref<User | null>(null)
const deleteTarget = ref<User | null>(null)
const isSaving = ref(false)
const isDeleting = ref(false)

const userForm = reactive({
  name: '',
  email: '',
  password: '',
  role: 'student' as UserRole,
  phone: ''
})

async function loadUsers() {
  isLoading.value = true
  error.value = null
  try {
    users.value = await adminService.getUsers({
      search: searchQuery.value || undefined,
      role: selectedRole.value || undefined
    })
  } catch (err: any) {
    error.value = err?.message || 'Gagal memuat pengguna'
  } finally {
    isLoading.value = false
  }
}

onMounted(loadUsers)

function openCreateModal() {
  editingUser.value = null
  Object.assign(userForm, {
    name: '',
    email: '',
    password: '',
    role: 'student',
    phone: ''
  })
  showModal.value = true
}

function openEditModal(u: User) {
  editingUser.value = u
  Object.assign(userForm, {
    name: u.name,
    email: u.email,
    password: '',
    role: u.role,
    phone: u.phone || ''
  })
  showModal.value = true
}

async function handleSaveUser() {
  if (!userForm.name.trim() || !userForm.email.trim()) {
    toast.error('Nama dan Email wajib diisi')
    return
  }

  isSaving.value = true
  try {
    if (editingUser.value) {
      const payload: any = {
        name: userForm.name,
        email: userForm.email,
        role: userForm.role,
        phone: userForm.phone
      }
      if (userForm.password) payload.password = userForm.password

      const updated = await adminService.updateUser(editingUser.value.id, payload)
      const idx = users.value.findIndex(u => u.id === editingUser.value!.id)
      if (idx !== -1) users.value[idx] = updated
      toast.success('Pengguna berhasil diperbarui!')
    } else {
      if (!userForm.password || userForm.password.length < 6) {
        toast.error('Password minimal 6 karakter')
        isSaving.value = false
        return
      }

      const created = await adminService.createUser(userForm)
      users.value.unshift(created)
      toast.success('Pengguna baru berhasil ditambahkan!')
    }
    showModal.value = false
  } catch (err: any) {
    toast.error(err?.message || 'Gagal menyimpan data pengguna')
  } finally {
    isSaving.value = false
  }
}

async function handleDeleteUser() {
  if (!deleteTarget.value) return
  isDeleting.value = true
  try {
    await adminService.deleteUser(deleteTarget.value.id)
    users.value = users.value.filter(u => u.id !== deleteTarget.value!.id)
    toast.success('Pengguna berhasil dihapus!')
    deleteTarget.value = null
  } catch (err: any) {
    toast.error(err?.message || 'Gagal menghapus pengguna')
  } finally {
    isDeleting.value = false
  }
}

const roleBadgeVariant: Record<UserRole, string> = {
  admin: 'danger',
  teacher: 'primary',
  student: 'default'
}
</script>

<template>
  <div class="p-4 md:p-6 max-w-6xl mx-auto space-y-6">
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
      <div>
        <h1 class="text-xl font-bold text-surface-900 dark:text-surface-100">Manajemen Pengguna</h1>
        <p class="text-xs text-surface-500 dark:text-surface-400 mt-0.5">Kelola akun Admin, Guru, dan Siswa pada sistem.</p>
      </div>

      <UiButton size="sm" @click="openCreateModal">
        <Plus class="w-4 h-4 mr-1" />
        Tambah Pengguna
      </UiButton>
    </div>

    <!-- Filter & Search Bar -->
    <div class="flex flex-col sm:flex-row items-center gap-3">
      <div class="w-full sm:flex-1">
        <UiInput
          v-model="searchQuery"
          placeholder="Cari nama atau email..."
          @input="loadUsers"
        />
      </div>

      <div class="w-full sm:w-48">
        <UiSelect
          v-model="selectedRole"
          :options="[
            { label: 'Semua Peran', value: '' },
            { label: 'Siswa', value: 'student' },
            { label: 'Guru', value: 'teacher' },
            { label: 'Admin', value: 'admin' }
          ]"
          @change="loadUsers"
        />
      </div>
    </div>

    <div v-if="isLoading">
      <UiSkeleton :rows="6" />
    </div>

    <UiErrorState v-else-if="error" :message="error" @retry="loadUsers" />

    <!-- Users Table -->
    <div v-else class="bg-white dark:bg-surface-900 rounded-2xl border border-surface-200 dark:border-surface-800 shadow-soft overflow-hidden">
      <div class="overflow-x-auto">
        <table class="w-full text-left text-sm">
          <thead class="bg-surface-50 dark:bg-surface-800/60 border-b border-surface-100 dark:border-surface-800 text-[11px] font-bold uppercase tracking-wider text-surface-500">
            <tr>
              <th class="py-3 px-4">Pengguna</th>
              <th class="py-3 px-4">Peran</th>
              <th class="py-3 px-4">Kontak</th>
              <th class="py-3 px-4">Terdaftar</th>
              <th class="py-3 px-4 text-right">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-surface-100 dark:divide-surface-800">
            <tr v-if="!users.length">
              <td colspan="5" class="py-8 text-center text-surface-400 text-xs">
                Tidak ada pengguna yang ditemukan.
              </td>
            </tr>
            <tr
              v-for="u in users"
              :key="u.id"
              class="hover:bg-surface-50/70 dark:hover:bg-surface-800/40 transition-colors"
            >
              <td class="py-3.5 px-4">
                <div class="flex items-center gap-3">
                  <UiAvatar :name="u.name" :src="u.avatar" size="sm" />
                  <div>
                    <p class="font-semibold text-surface-900 dark:text-surface-100 text-sm">{{ u.name }}</p>
                    <p class="text-xs text-surface-400">{{ u.email }}</p>
                  </div>
                </div>
              </td>
              <td class="py-3.5 px-4">
                <UiBadge :variant="roleBadgeVariant[u.role] as any" size="sm">
                  {{ u.role === 'admin' ? 'Admin' : u.role === 'teacher' ? 'Guru' : 'Siswa' }}
                </UiBadge>
              </td>
              <td class="py-3.5 px-4 text-xs text-surface-500">
                {{ u.phone || '-' }}
              </td>
              <td class="py-3.5 px-4 text-xs text-surface-400">
                {{ formatDate(u.created_at, { day: 'numeric', month: 'short', year: 'numeric' }) }}
              </td>
              <td class="py-3.5 px-4 text-right">
                <div class="flex items-center justify-end gap-1">
                  <button
                    type="button"
                    class="p-1.5 text-surface-400 hover:text-brand-600 dark:hover:text-brand-400 hover:bg-surface-100 dark:hover:bg-surface-800 rounded-lg transition-colors"
                    title="Edit Pengguna"
                    @click="openEditModal(u)"
                  >
                    <Pencil class="w-4 h-4" />
                  </button>
                  <button
                    type="button"
                    class="p-1.5 text-surface-400 hover:text-rose-600 dark:hover:text-rose-400 hover:bg-rose-50 dark:hover:bg-rose-950/30 rounded-lg transition-colors"
                    title="Hapus Pengguna"
                    @click="deleteTarget = u"
                  >
                    <Trash2 class="w-4 h-4" />
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Create / Edit User Modal -->
    <UiModal :show="showModal" :title="editingUser ? 'Edit Pengguna' : 'Tambah Pengguna Baru'" size="md" @close="showModal = false">
      <form @submit.prevent="handleSaveUser" class="space-y-4">
        <UiInput
          v-model="userForm.name"
          label="Nama Lengkap"
          placeholder="Nama pengguna"
          required
        />

        <UiInput
          v-model="userForm.email"
          type="email"
          label="Alamat Email"
          placeholder="user@sekolah.id"
          required
        />

        <UiInput
          v-model="userForm.password"
          type="password"
          :label="editingUser ? 'Password Baru (Kosongkan jika tidak ingin diubah)' : 'Password'"
          :placeholder="editingUser ? '••••••••' : 'Minimal 6 karakter'"
          :required="!editingUser"
        />

        <UiSelect
          v-model="userForm.role"
          label="Peran / Role"
          :options="[
            { label: 'Siswa (Student)', value: 'student' },
            { label: 'Guru (Teacher)', value: 'teacher' },
            { label: 'Admin Sistem', value: 'admin' }
          ]"
          required
        />

        <UiInput
          v-model="userForm.phone"
          label="Nomor Telepon (Opsional)"
          placeholder="08123456789"
        />
      </form>

      <template #footer>
        <div class="flex justify-end gap-2">
          <UiButton variant="outline" @click="showModal = false">Batal</UiButton>
          <UiButton :loading="isSaving" @click="handleSaveUser">
            {{ editingUser ? 'Simpan Perubahan' : 'Tambah Pengguna' }}
          </UiButton>
        </div>
      </template>
    </UiModal>

    <!-- Delete Confirm Modal -->
    <UiModal :show="!!deleteTarget" title="Hapus Pengguna?" size="sm" @close="deleteTarget = null">
      <p class="text-sm text-surface-600 dark:text-surface-400">
        Apakah Anda yakin ingin menghapus akun <strong>{{ deleteTarget?.name }}</strong> ({{ deleteTarget?.email }})?
      </p>
      <template #footer>
        <div class="flex justify-end gap-2">
          <UiButton variant="outline" @click="deleteTarget = null">Batal</UiButton>
          <UiButton variant="danger" :loading="isDeleting" @click="handleDeleteUser">Hapus</UiButton>
        </div>
      </template>
    </UiModal>
  </div>
</template>
