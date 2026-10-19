import { Action, getModule, Module, Mutation, VuexModule } from 'vuex-module-decorators'
import { createAPIKey, listAPIKeys, revokeAPIKey, rotateAPIKey } from '@/api/users'
import { APIKeyRecord, CreateAPIKeyRequest, CreatedAPIKeyResponse } from '@/models'
import store from '@/store'

export interface IAPIKeyState {
  keys: APIKeyRecord[]
}

@Module({ dynamic: true, store, name: 'apiKey' })
class APIKey extends VuexModule implements IAPIKeyState {
  public keys: APIKeyRecord[] = []

  @Mutation
  private SET_KEYS(keys: APIKeyRecord[]) {
    this.keys = keys
  }

  @Action
  public async Load() {
    const { data } = await listAPIKeys()
    this.SET_KEYS(data.list || [])
  }

  @Action
  public async Create(data: CreateAPIKeyRequest): Promise<CreatedAPIKeyResponse> {
    const response = await createAPIKey(data)
    await this.Load()
    return response.data as CreatedAPIKeyResponse
  }

  @Action
  public async Rotate(payload: { id: string, data: CreateAPIKeyRequest }): Promise<CreatedAPIKeyResponse> {
    const response = await rotateAPIKey(payload.id, payload.data)
    await this.Load()
    return response.data as CreatedAPIKeyResponse
  }

  @Action
  public async Revoke(id: string) {
    await revokeAPIKey(id)
    await this.Load()
  }
}

export const APIKeyModule = getModule(APIKey)
