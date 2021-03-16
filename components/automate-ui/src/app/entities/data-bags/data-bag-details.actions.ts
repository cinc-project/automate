import { HttpErrorResponse } from '@angular/common/http';
import { Action } from '@ngrx/store';
import { DataBagItems, DataBagItem } from './data-bags.model';

export enum DataBagItemsActionTypes {
  GET_ALL = 'DATA_BAG_ITEMS::GET_ALL',
  GET_ALL_SUCCESS = 'DATA_BAG_ITEMS::GET_ALL::SUCCESS',
  GET_ALL_FAILURE = 'DATA_BAG_ITEMS::GET_ALL::FAILURE',
  CREATE          = 'DATA_BAG_ITEMS::CREATE',
  CREATE_SUCCESS  = 'DATA_BAG_ITEMS::CREATE::SUCCESS',
  CREATE_FAILURE  = 'DATA_BAG_ITEMS::CREATE::FAILURE'
}

export interface DataBagItemsSuccessPayload {
  items: DataBagItems[];
  total: number;
}

export interface DataBagItemPayload {
  databagName: string;
  server_id: string;
  org_id: string;
  name: string;
  page: number;
  per_page: number;
}

export interface CreateDataBagItemPayload {
  name: string;
  id: string;
}

export class GetDataBagItems implements Action {
  readonly type = DataBagItemsActionTypes.GET_ALL;
  constructor(public payload: DataBagItemPayload) { }
}

export class GetDataBagItemsSuccess implements Action {
  readonly type = DataBagItemsActionTypes.GET_ALL_SUCCESS;
  constructor(public payload: DataBagItemsSuccessPayload) { }
}

export class GetDataBagItemsFailure implements Action {
  readonly type = DataBagItemsActionTypes.GET_ALL_FAILURE;
  constructor(public payload: HttpErrorResponse) { }
}

export class CreateDataBagItem implements Action {
  readonly type = DataBagItemsActionTypes.CREATE;
  constructor(public payload: { dataBagItem: DataBagItem }) { }
}

export class CreateDataBagItemSuccess implements Action {
  readonly type = DataBagItemsActionTypes.CREATE_SUCCESS;
  constructor(public payload: CreateDataBagItemPayload) { }
}

export class CreateDataBagItemFailure implements Action {
  readonly type = DataBagItemsActionTypes.CREATE_FAILURE;
  constructor(public payload: HttpErrorResponse) { }
}

export type DataBagItemsActions =
  | GetDataBagItems
  | GetDataBagItemsSuccess
  | GetDataBagItemsFailure
  | CreateDataBagItem
  | CreateDataBagItemSuccess
  | CreateDataBagItemFailure;
