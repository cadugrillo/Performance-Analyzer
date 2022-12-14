import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { HomeComponent } from './components/home/home.component';
import { ProfileComponent } from './components/profile/profile.component';
import { SignInComponent } from './components/sign-in/sign-in.component';
import { SignUpComponent } from './components/sign-up/sign-up.component';
import { AuthGuardService } from './services/auth-guard.service';
import { AnalyzeSignalsComponent } from './components/analyze-signals/analyze-signals.component';
import { ForgotPasswordComponent } from './components/forgot-password/forgot-password.component';
import { MqttClientComponent } from './components/mqtt-client/mqtt-client.component';

const routes: Routes = [

  { path: '', redirectTo: 'signIn', pathMatch: 'full' },
  { path: 'signIn',component: SignInComponent},
  { path: 'signUp',component: SignUpComponent},
  { path: 'forgotPassword',component: ForgotPasswordComponent},
  { path: 'profile',component: ProfileComponent, canActivate: [AuthGuardService]},
  { path: 'home', component: HomeComponent, canActivate: [AuthGuardService] },
  { path: 'analyze-signals', component: AnalyzeSignalsComponent,canActivate: [AuthGuardService] },
  { path: 'mqtt-client', component: MqttClientComponent,canActivate: [AuthGuardService] },
  { path: '**', redirectTo: 'signIn'},
 

];

@NgModule({
  imports: [RouterModule.forRoot(routes, {onSameUrlNavigation: 'reload'})],
  exports: [RouterModule]
})
export class AppRoutingModule { }
