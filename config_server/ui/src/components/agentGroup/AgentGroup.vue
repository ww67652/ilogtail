<template >
  <el-button type="primary" style="float:left" @click="addAgentGroup">新增</el-button>
  <el-dialog v-model="showCreateForm">
  <el-form  :model="createAgentGroup" :rules="createAgentGroupRules" ref="ruleFormRef">
    <el-form-item label="名字" prop="groupName">
      <el-input v-model="createAgentGroup.groupName" placeholder="名字"  />
    </el-form-item>
    <el-form-item label="描述">
      <el-input v-model="createAgentGroup.value" placeholder="描述" />
    </el-form-item>
    <el-form-item>
      <el-button type="primary" @click="onSubmit">保存</el-button>
    </el-form-item>
  </el-form>
  </el-dialog>

  <el-table :data="tableData">
    <el-table-column fixed prop="name" label="名称" min-width="200"/>
    <el-table-column prop="value" label="描述" width="560" >
      <template v-slot="scope">
        <div v-if="!scope.row.isEditing">
          {{ scope.row.value }}
        </div>
        <el-input v-else v-model="scope.row.value" placeholder="请输入描述" />
      </template>
    </el-table-column>
    <el-table-column label="采集Agent数" width="180" class-name="centered-column">
      <template v-slot="scope">
        <el-button type="primary" link @click="clickAgentCount(scope.row)">{{scope.row.agentList.length}}</el-button>
      </template>
    </el-table-column>
    <el-table-column  label="关联pipeline配置数" width="180" class-name="centered-column">
      <template #default="scope">
        <el-button type="primary" link @click="clickPipelineConfigCount(scope.row)">{{scope.row.appliedPipelineConfigList.length}}</el-button>
      </template>
    </el-table-column>
    <el-table-column label="关联instance配置数" width="180" class-name="centered-column">
      <template #default="scope">
        <el-button type="primary" link @click="clickInstanceConfigCount(scope.row)">{{scope.row.appliedInstanceConfigList.length}}</el-button>
      </template>
    </el-table-column>
    <el-table-column fixed="right" label="操作" min-width="120">
      <template #default="scope">
        <el-button link type="primary" size="small"
                   @click="scope.row.isEditing=!scope.row.isEditing">
          <div v-if="!scope.row.isEditing">
            编辑
          </div>
          <div v-else @click="saveValue(scope.row.name,scope.row.value)">
            保存
          </div>
        </el-button>
        <el-button link type="primary" size="small"
                   @click="deleteGroupByName(scope.row.name)"
                   :disabled="checkGroupNameIsDefault(scope.row.name)">
          删除
        </el-button>
      </template>
    </el-table-column>
  </el-table>


  <el-dialog v-model="showAgentInfoDialog" >
  <el-table :data="selectedRow.agentList">
    <el-table-column type="expand">
      <template #default="props">
          <el-table :data="props.row.agentConfigStatusList">
            <el-table-column label="type" prop="type"/>
            <el-table-column label="name" prop="name"/>
            <el-table-column label="status"  prop="status"/>
            <el-table-column label="message" prop="message"/>
          </el-table>
      </template>
    </el-table-column>
    <el-table-column property="instanceId" label="实例Id" width="300" />
    <el-table-column property="agentType" label="类型" width="200" />
    <el-table-column property="attributes" :formatter="formatAttributes" label="属性" width="100" />
    <el-table-column property="capabilities" label="能力" width="100"/>
    <el-table-column property="runningStatus" label="运行状态" />
    <el-table-column property="startupTime" label="开始时间" width="120"/>
  </el-table>
</el-dialog>

  <el-dialog v-model="showPipelineConfig" @close="closePipelineConfig">
    <el-table :data="selectedRow.allPipelineConfigList">
      <el-table-column property="name" label="名字" />
      <el-table-column property="version" label="版本"  />
      <el-table-column property="detail" label="详情"  />
      <el-table-column label="操作" min-width="100">
          <template v-slot="scope">
            <el-button link type="danger" size="small" v-if="scope.row.isApplied"
                       @click="cancelAppliedPipelineConfig(scope.row)">
              {{scope.row.description}}
            </el-button>

            <el-button link type="success" size="small" v-else
                       @click="appliedPipelineConfig(scope.row)">
              {{scope.row.description}}
            </el-button>
          </template>
      </el-table-column>
    </el-table>
  </el-dialog>

  <el-dialog v-model="showInstanceConfig" @close="closeInstanceConfig">
    <el-table :data="selectedRow.allInstanceConfigList">
      <el-table-column property="name" label="名字" />
      <el-table-column property="version" label="版本"  />
      <el-table-column property="detail" label="详情"  />
      <el-table-column label="操作" min-width="100">
        <template v-slot="scope">
          <el-button link type="danger" size="small" v-if="scope.row.isApplied"
                     @click="cancelAppliedInstanceConfig(scope.row)">
            {{scope.row.description}}
          </el-button>

          <el-button link type="success" size="small" v-else
                     @click="appliedInstanceConfig(scope.row)">
            {{scope.row.description}}
          </el-button>
        </template>
      </el-table-column>
    </el-table>
  </el-dialog>


</template>

<script >
import {
  createAgentGroup,
  deleteAgentGroup,
  getAppliedInstanceConfigsForAgentGroup,
  getAppliedPipelineConfigsForAgentGroup,
  listAgentGroups,
  listAgentsForAgentGroup, updateAgentGroup
} from "@/api/agentGroup";
import {messageShow} from "@/api/common"
import {
  applyPipelineConfigToAgentGroup,
  getPipelineConfig,
  listPipelineConfigs,
  removePipelineConfigFromAgentGroup
} from "@/api/pipelineConfig";
import {
  applyInstanceConfigToAgentGroup, getInstanceConfig,
  listInstanceConfigs,
  removeInstanceConfigFromAgentGroup
} from "@/api/instanceConfig";
import {getInstanceConfigStatusList, getPipelineConfigStatusList} from "@/api/agent";

export default {
  name: "AgentGroup",
  data() {
    return {
      tableData: [],
      selectedRow: null,

      showAgentInfoDialog: false,
      showCreateForm: false,
      showPipelineConfig: false,
      showInstanceConfig: false,

      createAgentGroup: {
        groupName: "",
        value: "",
      },
      createAgentGroupRules: {
        groupName: [
          {required: true, message: '请输入用户名', trigger: 'blur'},
        ]
      },

    }
  },
  async created() {
    await this.initAllTable()
    this.initDataMap()
  },

  methods: {

    addAgentGroup() {
      this.createAgentGroup={}
      this.showCreateForm = true
    },

    async onSubmit() {
      this.$refs["ruleFormRef"].validate(async valid => {
        if (valid) {
          let data = await createAgentGroup(this.createAgentGroup.groupName, this.createAgentGroup.value)
          messageShow(data, "新增成功")
          this.showCreateForm = false
          await this.initAllTable()
        }
      })

    },

    formatAttributes(row, column, cellValue) {
      if (typeof cellValue === 'object') {
        return Object.entries(cellValue)
            .map(([key, value]) => `${key}: ${value}`)
            .join('\n\n');
      }
      return cellValue;
    },

    async initAllTable() {
      let data = await listAgentGroups()
      if (!data.commonResponse.status) {
        this.tableData = await Promise.all(data.agentGroupsList.map(async item => {
          console.log((await getAppliedInstanceConfigsForAgentGroup(item.name)))
          return {
            name: item.name,
            value: item.value,
            agentList: (await listAgentsForAgentGroup(item.name)).agentsList,
            appliedPipelineConfigList: (await getAppliedPipelineConfigsForAgentGroup(item.name)).configNamesList.map(res => {
              return {name: res}
            }),
            appliedInstanceConfigList: (await getAppliedInstanceConfigsForAgentGroup(item.name)).configNamesList.map(res => {
              return {name: res}
            })
          }
        }))

        if (this.tableData.length>1){
          let temp = this.tableData[0]
          for(let i=0;i<this.tableData.length;i++){
            if(this.tableData[i].name==="default"){
              this.tableData[0]=this.tableData[i]
              this.tableData[i]=temp
              break
            }
          }
        }
      }
    },
    initDataMap() {
      this.tableData.forEach(item => {
        item.isEditing=false
      })
    },

    findSelectedRow(name){
      for(let row of this.tableData){
        if(row.name===name){
          return row
        }
      }
      return null
    },

    async clickAgentCount(row) {
      this.selectedRow = this.findSelectedRow(row.name)
      const agentCount = this.selectedRow.agentList.length
      this.showAgentInfoDialog = agentCount !== 0
      if (this.showAgentInfoDialog){
        for(let agent of this.selectedRow.agentList){
          agent.agentConfigStatusList=[]
          let pipelineConfigStatusList=(await getPipelineConfigStatusList(agent.instanceId)).agentConfigStatusList
          for(let pipelineConfigStatus of pipelineConfigStatusList){
            pipelineConfigStatus.type="pipelineConfig"
          }

          let instanceConfigStatusList=(await getInstanceConfigStatusList(agent.instanceId)).agentConfigStatusList
          for(let instanceConfigStatus of instanceConfigStatusList){
            instanceConfigStatus.type="instanceConfig"
          }
          agent.agentConfigStatusList=agent.agentConfigStatusList.concat(pipelineConfigStatusList)
          agent.agentConfigStatusList=agent.agentConfigStatusList.concat(instanceConfigStatusList)
        }
        console.log(this.selectedRow.agentList)
      }
    },

    async saveValue(name, value) {
      const data = await updateAgentGroup(name, value)
      messageShow(data, "保存成功")
    },

    checkGroupNameIsDefault(name) {
      return name === "default"
    },

    async deleteGroupByName( name) {
      let data = await deleteAgentGroup(name)
      messageShow(data, "删除成功")
      await this.initAllTable()
    },

    // ______________________________
    async clickPipelineConfigCount(row) {
      this.selectedRow = this.findSelectedRow(row.name)
      this.showPipelineConfig = true
      await this.getPipelineConfigInfo(this.selectedRow)
    },

    async getPipelineConfigInfo(row) {
      let allPipeConfigList = (await listPipelineConfigs()).configDetailsList
      let appliedPipelineConfigList = await Promise.all(
          row.appliedPipelineConfigList.map(async config => (await getPipelineConfig(config.name)).configDetail))

      for (let config of allPipeConfigList) {
        config.description = "应用"
        config.isApplied = false
        for (let appliedConfig of appliedPipelineConfigList) {
          if (config.name === appliedConfig.name) {
            config.isApplied = true
            config.description = "取消应用"
            break
          }
        }
      }
      allPipeConfigList.sort((c1, c2) => {
        return c2.isApplied - c1.isApplied
      })
      row.allPipelineConfigList = allPipeConfigList
    },


    async cancelAppliedPipelineConfig(row) {
      let data = await removePipelineConfigFromAgentGroup(row.name, this.selectedRow.name)
      messageShow(data, "取消应用配置成功")
      if (!data.commonResponse.status) {
        row.description = "应用"
        row.isApplied = !row.isApplied
      }
    },

    async appliedPipelineConfig(row) {
      let data = await applyPipelineConfigToAgentGroup(row.name, this.selectedRow.name)
      messageShow(data, "应用配置成功")
      if (!data.commonResponse.status) {
        row.description = "取消应用"
        row.isApplied = !row.isApplied
      }
    },

    async closePipelineConfig() {
      this.showPipelineConfig=false
      await this.initAllTable()
    },


  // ____________________________


  // ______________________________
  async clickInstanceConfigCount(row) {
    this.selectedRow = this.findSelectedRow(row.name)
    this.showInstanceConfig = true
    await this.getInstanceConfigInfo(this.selectedRow)
  },

  async getInstanceConfigInfo(row) {
    let allInstanceConfigList = (await listInstanceConfigs()).configDetailsList
    console.log(allInstanceConfigList)
    let appliedInstanceConfigList = await Promise.all(
        row.appliedInstanceConfigList.map(async config => (await getInstanceConfig(config.name)).configDetail))

    for (let config of allInstanceConfigList) {
      config.description = "应用"
      config.isApplied = false
      for (let appliedConfig of appliedInstanceConfigList) {
        if (config.name === appliedConfig.name) {
          config.isApplied = true
          config.description = "取消应用"
          break
        }
      }
    }
    allInstanceConfigList.sort((c1, c2) => {
      return c2.isApplied - c1.isApplied
    })
    row.allInstanceConfigList = allInstanceConfigList
  },


  async cancelAppliedInstanceConfig(row) {
    let data = await removeInstanceConfigFromAgentGroup(row.name, this.selectedRow.name)
    messageShow(data, "取消应用配置成功")
    if (!data.commonResponse.status) {
      row.description = "应用"
      row.isApplied = !row.isApplied
    }
  },

  async appliedInstanceConfig(row) {
    let data = await applyInstanceConfigToAgentGroup(row.name, this.selectedRow.name)
    messageShow(data, "应用配置成功")
    if (!data.commonResponse.status) {
      row.description = "取消应用"
      row.isApplied = !row.isApplied
    }
  },

  async closeInstanceConfig() {
    this.showInstanceConfig=false
    await this.initAllTable()
  },

},

// ____________________________

}
</script>

<style>

</style>
