import { Toaster } from "@/components/ui/toaster";
import { Toaster as Sonner } from "@/components/ui/sonner";
import { TooltipProvider } from "@/components/ui/tooltip";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import Index from "./pages/Index";
import Dashboard from "./pages/Dashboard";
import LogViewer from "./pages/LiveLogging";
import Documentation from "./pages/Documentation";
import NotFound from "./pages/NotFound";
import { inject } from "@vercel/analytics";
import { useEffect } from "react";
import {SpeedInsights} from "@vercel/speed-insights/react";


const queryClient = new QueryClient();

const App = () => {
  useEffect(() => {
    inject();
  }, []);
  return (
    <QueryClientProvider client={queryClient}>
      <TooltipProvider>
        <Toaster />
        <Sonner position="top-right" theme="dark" />
        <BrowserRouter>
          <Routes>
            <Route path="/" element={<Index />} />
            <Route path="/docs" element={<Documentation />} />
            <Route path="/demo-logs/:streamId" element={<Dashboard />} />
            <Route path="/live-logs/:streamId" element={<LogViewer />}></Route>
            <Route path="*" element={<NotFound />} />
          </Routes>
        </BrowserRouter>
        <SpeedInsights/> {/* For Vercel Speed Insights */}
      </TooltipProvider>
    </QueryClientProvider>
  )
};

export default App;
