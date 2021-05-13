import { EntityState, EntityAdapter, createEntityAdapter } from '@ngrx/entity';
import { set, pipe } from 'lodash/fp';

import { EntityStatus } from 'app/entities/entities';
import { NodeActionTypes, NodeActions } from './infra-nodes.actions';
import { InfraNode } from './infra-nodes.model';

export interface InfraNodeEntityState extends EntityState<InfraNode> {
  nodesStatus: EntityStatus;
  getAllStatus: EntityStatus;
  getStatus: EntityStatus;
  updateEnvStatus: EntityStatus;
  updateTagsStatus: EntityStatus;
  nodeList: {
    items: InfraNode[],
    total: number
  };
  nodeTags: string[];
  nodeEnvironment: string;
}

const GET_ALL_STATUS = 'getAllStatus';
const GET_STATUS = 'getStatus';
const UPDATE_ENVIRONMENT_STATUS = 'updateEnvStatus';
const UPDATE_TAGS_STATUS = 'updateTagsStatus';

export const nodeEntityAdapter: EntityAdapter<InfraNode> = createEntityAdapter<InfraNode>({
  selectId: (infraNode: InfraNode) => infraNode.name
});

export const InfraNodeEntityInitialState: InfraNodeEntityState =
  nodeEntityAdapter.getInitialState(<InfraNodeEntityState>{
    getAllStatus: EntityStatus.notLoaded
  });

export function infraNodeEntityReducer(
  state: InfraNodeEntityState = InfraNodeEntityInitialState,
  action: NodeActions): InfraNodeEntityState {

  switch (action.type) {
    case NodeActionTypes.GET_ALL:
      return set(GET_ALL_STATUS, EntityStatus.loading, nodeEntityAdapter.removeAll(state));

    case NodeActionTypes.GET_ALL_SUCCESS:
      return pipe(
        set(GET_ALL_STATUS, EntityStatus.loadingSuccess),
        set('nodeList.items', action.payload.nodes || []),
        set('nodeList.total', action.payload.total || 0)
        )(state) as InfraNodeEntityState;

    case NodeActionTypes.GET_ALL_FAILURE:
      return set(GET_ALL_STATUS, EntityStatus.loadingFailure, state);

    case NodeActionTypes.GET:
      return set(
        GET_STATUS,
        EntityStatus.loading,
        nodeEntityAdapter.removeAll(state)
      ) as InfraNodeEntityState;

    case NodeActionTypes.GET_SUCCESS:
      return set(GET_STATUS, EntityStatus.loadingSuccess,
        nodeEntityAdapter.addOne(action.payload, state));

    case NodeActionTypes.GET_FAILURE:
      return set(GET_STATUS, EntityStatus.loadingFailure, state);

    case NodeActionTypes.UPDATE_ENVIRONMENT:
      return set(UPDATE_ENVIRONMENT_STATUS, EntityStatus.loading, state);

    case NodeActionTypes.UPDATE_ENVIRONMENT_SUCCESS:
      return pipe(
        set(UPDATE_ENVIRONMENT_STATUS, EntityStatus.loadingSuccess),
        set('nodeEnvironment', action.payload.environment || [])
        )(state) as InfraNodeEntityState;

    case NodeActionTypes.UPDATE_ENVIRONMENT_FAILURE:
      return set(UPDATE_ENVIRONMENT_STATUS, EntityStatus.loadingFailure, state);

    case NodeActionTypes.UPDATE_TAGS:
      return set(UPDATE_TAGS_STATUS, EntityStatus.loading, state);

    case NodeActionTypes.UPDATE_TAGS_SUCCESS:
      return pipe(
        set(UPDATE_TAGS_STATUS, EntityStatus.loadingSuccess),
        set('nodeTags', action.payload.tags || [])
        )(state) as InfraNodeEntityState;

    case NodeActionTypes.UPDATE_TAGS_FAILURE:
      return set(UPDATE_TAGS_STATUS, EntityStatus.loadingFailure, state);

    default:
      return state;
  }
}

export const getEntityById = (id: string) =>
  (state: InfraNodeEntityState) => state.entities[id];
