<template>
  <div class="m-4 p-2">
    <h1 class="text-center text-amber-500">Override-Gui</h1>
    <div class="flex">
      <div class="mr-4 min-w-50">
        <Card>
          <template #title>
            <h2 class="m-1">服务管理</h2>
          </template>
          <template #content>
            <div class="mt-4 flex flex-col">
              <Button :label="serverButtonLabel" @click="toggleServer" :severity="serverButtonSeverity" class="mt-4" />
              <Button label="复制 vscode 配置" severity="info" class="mt-4" @click="copyVscodeConfig" />
              <Button label="连通性测试" severity="contrast" class="mt-4" @click="testConnection" />
            </div>
          </template>
        </Card>
      </div>
      <div>
        <Card>
          <template #title>
            <h2 class="m-1">配置管理</h2>
          </template>
          <template #content>
            <div class="flex flex-wrap p-4">
              <div class="flex items-center mb-4 w-full" v-for="(value, key) in form" :key="key">
                <label :for="key" class="mr-4 min-w-50 text-right">{{ key }}</label>
                <InputText v-if="!isNumber(key)" v-model="form[key]" :id="key" class="w-full" />
                <InputNumber v-if="isNumber(key)" v-model="form[key]" :id="key" class="w-full" />
              </div>
            </div>
            <div class="flex justify-center mt-4">
              <Button label="更新配置并重启" class="p-button-primary" @click="updateConfig" />
            </div>
          </template>
        </Card>
      </div>
    </div>
    <Toast />
  </div>
</template>

<script>
import { ref, watch, onMounted } from 'vue';
import { useToast } from 'primevue/usetoast';
import { BackendService } from "../bindings/override-gui/";

export default {
  setup() {
    const toast = useToast();
    const serverRunning = ref(false);
    const serverButtonLabel = ref('启动');
    const serverButtonSeverity = ref('success');

    const form = ref({
      bind: '127.0.0.1:8181',
      timeout: 600,
      api_key: '',
      codex_api_base: 'https://api.deepseek.com/beta/v1',
      codex_max_tokens: 500,
      code_instruct_model: 'deepseek-coder',
      chat_api_base: 'https://api.deepseek.com/v1',
      chat_max_tokens: 4096,
      chat_model_default: 'deepseek-chat',
      chat_locale: 'zh_CN'
    });

    const isNumber = (key) => {
      const numberFields = ['timeout', 'codex_max_tokens', 'chat_max_tokens'];
      return numberFields.includes(key);
    };

    const toggleServer = async () => {
      if (serverRunning.value) {
        const res = await BackendService.StopServer();
        if (res.status === 'success') {
          serverRunning.value = false;
          toast.add({ severity: 'success', summary: '成功', detail: res.msg, life: 3000 });
        } else {
          toast.add({ severity: 'error', summary: '失败', detail: res.msg, life: 3000 });
        }
      } else {
        const res = await BackendService.StartServer();
        if (res.status === 'success') {
          serverRunning.value = true;
          toast.add({ severity: 'success', summary: '成功', detail: res.msg, life: 3000 });
        } else {
          toast.add({ severity: 'error', summary: '失败', detail: res.msg, life: 3000 });
        }
      }
    };

    const updateConfig = async () => {
      const payload = {
        bind: form.value.bind,
        proxy_url: "",
        timeout: form.value.timeout,
        codex_api_base: form.value.codex_api_base,
        codex_api_key: form.value.api_key,
        codex_api_organization: "",
        codex_api_project: "",
        codex_max_tokens: form.value.codex_max_tokens,
        code_instruct_model: form.value.code_instruct_model,
        chat_api_base: form.value.chat_api_base,
        chat_api_key: form.value.api_key,
        chat_api_organization: "",
        chat_api_project: "",
        chat_max_tokens: form.value.chat_max_tokens,
        chat_model_default: form.value.chat_model_default,
        chat_model_map: {},
        chat_locale: form.value.chat_locale,
        auth_token: ""
      };
      const res = await BackendService.UpdateConfig(JSON.stringify(payload));
      if (res.status === 'success') {
        toast.add({ severity: 'success', summary: '成功', detail: res.msg, life: 3000 });
        if (serverRunning.value) {
          await BackendService.StopServer();
          const startRes = await BackendService.StartServer();
          if (startRes.status === 'success') {
            await testConnection();
            serverRunning.value = true;
          } else {
            toast.add({ severity: 'error', summary: '失败', detail: startRes.msg, life: 3000 });
          }
        } else {
          const startRes = await BackendService.StartServer();
          if (startRes.status === 'success') {
            await testConnection();
            serverRunning.value = true;
          } else {
            toast.add({ severity: 'error', summary: '失败', detail: startRes.msg, life: 3000 });
          }
        }
      } else {
        toast.add({ severity: 'error', summary: '失败', detail: res.msg, life: 3000 });
      }
    };

    const testConnection = async () => {
      try {
        const response = await fetch(`http://${form.value.bind}/_ping`, { mode: 'cors' });
        if (response.ok) {
          toast.add({ severity: 'success', summary: '成功', detail: '服务器已启动', life: 3000 });
        } else {
          toast.add({ severity: 'error', summary: '失败', detail: '服务器未启动', life: 3000 });
        }
      } catch (error) {
        toast.add({ severity: 'error', summary: '失败', detail: '服务器未启动'});
      }
    };

    const copyVscodeConfig = async () => {
      let config = {
        "github.copilot.advanced": {
          "debug.overrideCAPIUrl": `http://${form.value.bind}/v1`,
          "debug.overrideProxyUrl": `http://${form.value.bind}`,
          "debug.chatOverrideProxyUrl": `http://${form.value.bind}/v1/chat/completions`,
          "authProvider": "github-enterprise"
        },
        "github-enterprise.uri": "https://cocopilot.org",
      };

      let configStr = JSON.stringify(config, null, 2);
      configStr = configStr.slice(1, -1);

      await navigator.clipboard.writeText(`${configStr.trimEnd()},\n`);
      toast.add({ severity: 'success', summary: '成功', detail: 'vscode 配置已复制到剪贴板', life: 3000 });
    };

    const readConfig = async () => {
      const res = await BackendService.ReadConfig();
      if (res.status === "success") {
        const config = res.data;
        form.value.bind = config.bind;
        form.value.timeout = config.timeout;
        form.value.api_key = config.chat_api_key;
        form.value.codex_api_base = config.codex_api_base;
        form.value.codex_max_tokens = config.codex_max_tokens;
        form.value.code_instruct_model = config.code_instruct_model;
        form.value.chat_api_base = config.chat_api_base;
        form.value.chat_max_tokens = config.chat_max_tokens;
        form.value.chat_model_default = config.chat_model_default;
        form.value.chat_locale = config.chat_locale;
      } else {
        toast.add({ severity: 'error', summary: '失败', detail: res.msg, life: 3000 });
      }
    }

    watch(serverRunning, (newVal) => {
      serverButtonLabel.value = newVal ? '停止' : '启动';
      serverButtonSeverity.value = newVal ? 'danger' : 'success';
    });

    onMounted(async () => {
      await readConfig();
    });

    return {
      form,
      updateConfig,
      toggleServer,
      serverRunning,
      isNumber,
      copyVscodeConfig,
      testConnection,
      readConfig,
      serverButtonLabel,
      serverButtonSeverity,
    };
  },
};
</script>

<style>
html {
  padding: 0;
  margin: 0;
  background-color: #F8FAFC;
}
</style>
