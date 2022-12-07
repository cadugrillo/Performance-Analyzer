import { Injectable } from '@angular/core';
import { IMqttMessage, MqttService } from "ngx-mqtt";
import { Observable } from "rxjs";

@Injectable({
  providedIn: 'root'
})
export class MqttClientService {

  constructor(private mqttClientService: MqttService) { }

  topic(topic: string, qos: 0 | 1 | 2): Observable<IMqttMessage> {
    return this.mqttClientService.observe(topic, {qos: qos});
  }
}
