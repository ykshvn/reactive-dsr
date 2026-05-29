import { HttpClient, HttpParams } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';

export interface Node {
  id: number;
  x: number;
  y: number;
  neighbors: number[];
}

export interface Edge {
  from: number;
  to: number;
}

export interface GraphResponse {
  nodes: Node[];
  edges: Edge[];
  node_count: number;
}

export interface StepResponse {
  status: string;
  step: number;
  queueLength: number;
  event?: Event;
  message?: string;
}

export interface Event {
  type: string;
  step: number;
  payload: any;
  timestamp: number;
}

@Injectable({
  providedIn: 'root',
})
export class GraphService {
  private baseUrl = 'http://localhost:6969/api';
  constructor(private http: HttpClient) {}

  generateGraph(nodesCount: number = 20): Observable<GraphResponse> {
    return this.http.get<GraphResponse>(
      `${this.baseUrl}/graph/generate?nodes=${nodesCount}`,
    );
  }

  startRouteDiscovery(src: number, dst: number): Observable<any> {
    const params = new HttpParams()
      .set('src', src.toString())
      .set('dst', dst.toString());
    return this.http.get(`${this.baseUrl}/simulation/start`, { params });
  }

  nextStep(): Observable<StepResponse> {
    return this.http.get<StepResponse>(`${this.baseUrl}/simulation/step`);
  }

  runSimulation(source: number, destination: number): Observable<any> {
    const params = new HttpParams()
      .set('src', source.toString())
      .set('dst', destination.toString());

    return this.http.get(`${this.baseUrl}/simulation/run`, { params });
  }
}
