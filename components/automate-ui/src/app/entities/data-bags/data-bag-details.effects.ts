import { Injectable } from '@angular/core';
import { HttpErrorResponse } from '@angular/common/http';
import { Actions, Effect, ofType } from '@ngrx/effects';
import { of as observableOf } from 'rxjs';
import { catchError, mergeMap, map, filter } from 'rxjs/operators';
import { CreateNotification } from 'app/entities/notifications/notification.actions';
import { Type } from 'app/entities/notifications/notification.model';
import { HttpStatus } from 'app/types/types';

import {
  GetDataBagItems,
  GetDataBagItemsSuccess,
  GetDataBagItemsFailure,
  DataBagItemsActionTypes,
  DataBagItemsSuccessPayload,
  CreateDataBagItem,
  CreateDataBagItemSuccess,
  CreateDataBagItemPayload,
  CreateDataBagItemFailure
} from './data-bag-details.actions';

import { DataBagsRequests } from './data-bags.requests';

@Injectable()
export class DataBagItemsEffects {
  constructor(
    private actions$: Actions,
    private requests: DataBagsRequests
  ) { }

  @Effect()
  getDataBagItems$ = this.actions$.pipe(
    ofType(DataBagItemsActionTypes.GET_ALL),
    mergeMap(( action: GetDataBagItems) =>
      this.requests.getDataBagItems(action.payload).pipe(
        map((resp: DataBagItemsSuccessPayload) => new GetDataBagItemsSuccess(resp)),
        catchError((error: HttpErrorResponse) =>
          observableOf(new GetDataBagItemsFailure(error))))));

  @Effect()
  getDataBagItemsFailure$ = this.actions$.pipe(
    ofType(DataBagItemsActionTypes.GET_ALL_FAILURE),
    map(({ payload }: GetDataBagItemsFailure) => {
      const msg = payload.error.error;
      return new CreateNotification({
        type: Type.error,
        message: `Could not get infra data bag items: ${msg || payload.error}`
      });
    }));

    @Effect()
  createDataBagItem$ = this.actions$.pipe(
    ofType(DataBagItemsActionTypes.CREATE),
    mergeMap(({ payload: { dataBagItem } }: CreateDataBagItem) =>
      this.requests.createDataBagItem(dataBagItem).pipe(
        map((resp: CreateDataBagItemPayload) => new CreateDataBagItemSuccess(resp)),
        catchError((error: HttpErrorResponse) =>
          observableOf(new CreateDataBagItemFailure(error))))));

  @Effect()
  createDataBagItemSuccess$ = this.actions$.pipe(
    ofType(DataBagItemsActionTypes.CREATE_SUCCESS),
    map(({ payload: { name : name } }: CreateDataBagItemSuccess) => {
      return new CreateNotification({
        type: Type.info,
        message: `Successfully Created Data Bag ${name}.`
      });
    }));

  @Effect()
  createDataBagItemFailure$ = this.actions$.pipe(
    ofType(DataBagItemsActionTypes.CREATE_FAILURE),
    filter(({ payload }: CreateDataBagItemFailure) => payload.status !== HttpStatus.CONFLICT),
    map(({ payload }: CreateDataBagItemFailure) => {
      return new CreateNotification({
        type: Type.error,
        message: `Could Not Create Data Bag: ${payload.error.error || payload}.`
      });
    }));
}
