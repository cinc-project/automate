import { ComponentFixture, TestBed } from '@angular/core/testing';

import { DataFeedConfigDetailsComponent } from './data-feed-config-details.component';

describe('DataFeedConfigDetailsComponent', () => {
  let component: DataFeedConfigDetailsComponent;
  let fixture: ComponentFixture<DataFeedConfigDetailsComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ DataFeedConfigDetailsComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(DataFeedConfigDetailsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
