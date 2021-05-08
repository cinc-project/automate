import { Injectable } from '@angular/core';
import { HttpErrorResponse } from '@angular/common/http';
import { Actions, createEffect, ofType } from '@ngrx/effects';
import { of as observableOf } from 'rxjs';
import { catchError, mergeMap, map, filter } from 'rxjs/operators';
import { HttpStatus } from 'app/types/types';
import { CreateNotification } from 'app/entities/notifications/notification.actions';
import { Type } from 'app/entities/notifications/notification.model';

import {
  GetNodes,
  GetNode,
  GetNodesSuccess,
  GetNodeSuccess,
  NodesSuccessPayload,
  GetNodesFailure,
  GetNodeFailure,
  UpdateNodeEnvironment,
  UpdateNodeEnvironmentSuccess,
  UpdateNodeEnvironmentFailure,
  NodeActionTypes
} from './infra-nodes.actions';

import {
  InfraNodeRequests
} from './infra-nodes.requests';

@Injectable()
export class InfraNodeEffects {
  constructor(
    private actions$: Actions,
    private requests: InfraNodeRequests
  ) { }

  getNodes$ = createEffect(() =>
    this.actions$.pipe(
    ofType(NodeActionTypes.GET_ALL),
    mergeMap((action: GetNodes) =>
      this.requests.getNodes(action.payload).pipe(
        map((resp: NodesSuccessPayload) => new GetNodesSuccess(resp)),
        catchError((error: HttpErrorResponse) =>
          observableOf(new GetNodesFailure(error)))))));

  getNodesFailure$ = createEffect(() =>
    this.actions$.pipe(
    ofType(NodeActionTypes.GET_ALL_FAILURE),
    map(({ payload }: GetNodesFailure) => {
      const msg = payload.error.error;
      return new CreateNotification({
        type: Type.error,
        message: `Could not get nodes: ${msg || payload.error}`
      });
    })));

  getNode$ = createEffect(() =>
    this.actions$.pipe(
      ofType(NodeActionTypes.GET),
      mergeMap(({ payload: { server_id, org_id, name } }: GetNode) =>
        this.requests.getNode(server_id, org_id, name).pipe(
          map((resp) => new GetNodeSuccess(resp)),
          catchError((error: HttpErrorResponse) => observableOf(new GetNodeFailure(error)))))));

  getNodeFailure$ = createEffect(() =>
    this.actions$.pipe(
      ofType(NodeActionTypes.GET_FAILURE),
      map(({ payload }: GetNodeFailure) => {
        const msg = payload.error.error;
        return new CreateNotification({
          type: Type.error,
          message: `Could not get node: ${msg || payload.error}`
        });
    })));
  
  updateNodeEnvironment$ = createEffect(() =>
    this.actions$.pipe(
      ofType(NodeActionTypes.UPDATE),
      mergeMap(({ payload }: UpdateNodeEnvironment) =>
        this.requests.updateNodeEnvironment(payload.node).pipe(
          map((resp) => new UpdateNodeEnvironmentSuccess(resp)),
          catchError((error: HttpErrorResponse) =>
            observableOf(new UpdateNodeEnvironmentFailure(error)))))));

  updateNodeEnvironmentSuccess$ = createEffect(() =>
    this.actions$.pipe(
      ofType(NodeActionTypes.UPDATE_SUCCESS),
      map(({ }: UpdateNodeEnvironmentSuccess) => new CreateNotification({
        type: Type.info,
        message: `Successfully updated node environment.`
      }))));

  updateNodeEnvironmentFailure$ = createEffect(() =>
    this.actions$.pipe(
      ofType(NodeActionTypes.UPDATE_FAILURE),
      filter(({ payload }: UpdateNodeEnvironmentFailure) => payload.status !== HttpStatus.CONFLICT),
      map(({ payload }: UpdateNodeEnvironmentFailure) => new CreateNotification({
        type: Type.error,
        message: `Could not update node environment: ${payload.error.error || payload}.`
      }))));
}
