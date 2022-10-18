import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { HomeComponent } from './components/home/home.component';
import { AuthGuardService } from './services/auth-guard.service';
import { AnalyzeSignalsComponent } from './components/analyze-signals/analyze-signals.component';
import { MqttClientComponent } from './components/mqtt-client/mqtt-client.component';

const routes: Routes = [

  { path: '', redirectTo: 'home', pathMatch: 'full' },
  { path: 'home', component: HomeComponent, canActivate: [AuthGuardService] },
  { path: 'analyze-signals', component: AnalyzeSignalsComponent,canActivate: [AuthGuardService] },
  { path: 'mqtt-client', component: MqttClientComponent,canActivate: [AuthGuardService] },
  { path: '**', redirectTo: 'home'},
 

];

@NgModule({
  imports: [RouterModule.forRoot(routes, {onSameUrlNavigation: 'reload'})],
  exports: [RouterModule]
})
export class AppRoutingModule { }
