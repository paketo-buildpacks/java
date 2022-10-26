package io.paketo;

import jakarta.servlet.annotation.WebServlet;
import jakarta.servlet.http.HttpServlet;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;

import org.apache.catalina.util.ServerInfo;

import java.io.IOException;

@WebServlet("/tomcat")
public class TomcatInfoServlet extends HttpServlet {

    @Override
    protected void doGet(HttpServletRequest req, HttpServletResponse resp) throws IOException {
        resp.getOutputStream().println(ServerInfo.getServerInfo());
    }
}
