import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';

import { AppRoutingModule } from './app-routing.module';
import { HttpClientModule, HTTP_INTERCEPTORS } from '@angular/common/http';
import {MatInputModule} from '@angular/material/input';
import {MatButtonModule} from '@angular/material/button';
import {MatIconModule} from '@angular/material/icon';
import {MatFormFieldModule} from '@angular/material/form-field';
import {MatCheckboxModule} from '@angular/material/checkbox';
import {MatCardModule} from '@angular/material/card';
import {MatDialogModule} from '@angular/material/dialog';
import {MatToolbarModule} from '@angular/material/toolbar';
import {MatTabsModule} from '@angular/material/tabs';
import {MatMenuModule} from '@angular/material/menu';
import {MatListModule} from '@angular/material/list';
import {MatSelectModule} from '@angular/material/select';
import { MessagePopupComponent } from './components/message-popup/message-popup.component';
import { WaitPopupComponent } from './components/wait-popup/wait-popup.component';
import {MatProgressSpinnerModule} from '@angular/material/progress-spinner';
import {MatProgressBarModule} from '@angular/material/progress-bar';
import {MatTableModule} from '@angular/material/table';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import {MatTooltipModule} from '@angular/material/tooltip';

import { AppComponent } from './app.component';
import { HomeComponent } from './components/home/home.component';
import { ProfileComponent } from './components/profile/profile.component';
import { SignInComponent } from './components/sign-in/sign-in.component';
import { SignUpComponent } from './components/sign-up/sign-up.component';
import { FormsModule } from '@angular/forms';
import { TokenInterceptor } from './token.interceptor';
import { MatPaginatorModule } from '@angular/material/paginator';
import { SignalsService } from "./services/signals.service";
import { AnalyzeSignalsComponent } from './components/analyze-signals/analyze-signals.component';


@NgModule({
  declarations: [
    AppComponent,
    HomeComponent,
    ProfileComponent,
    SignInComponent,
    SignUpComponent,
    MessagePopupComponent,
    WaitPopupComponent,
    AnalyzeSignalsComponent
  ],
  imports: [
    AppRoutingModule,
    BrowserModule,
    FormsModule,
    HttpClientModule,
    MatFormFieldModule,
    MatInputModule,
    MatButtonModule,
    MatIconModule,
    MatCheckboxModule,
    MatDialogModule,
    MatToolbarModule,
    MatTabsModule,
    MatListModule,
    MatCardModule,
    MatMenuModule,
    MatSelectModule,
    MatProgressSpinnerModule,
    MatProgressBarModule,
    MatTableModule,
    MatPaginatorModule,
    BrowserAnimationsModule,
    MatTooltipModule
  ],
  providers: [SignalsService, {
    provide: HTTP_INTERCEPTORS,
    useClass: TokenInterceptor,
    multi: true
  },],
  bootstrap: [AppComponent]
})
export class AppModule { }
